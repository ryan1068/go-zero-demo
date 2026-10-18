package logic

import (
	"context"
	"database/sql"
	"errors"

	"go-zero-demo/order-rpc/internal/svc"
	"go-zero-demo/order-rpc/pb/order"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartLogic {
	return &AddCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddCartLogic) AddCart(in *order.AddCartReq) (*order.CartItemResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.CartItemResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.ProductId <= 0 || in.Quantity < 1 || in.Quantity > maxQuantity {
		return &order.CartItemResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin cart transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	var (
		productName string
		cover       string
		price       int64
		stock       int32
		status      int32
	)
	err = tx.QueryRowContext(l.ctx,
		"SELECT name, COALESCE(cover_image, ''), price, stock, status FROM looklook_product.product WHERE id = ? AND deleted_at IS NULL FOR UPDATE",
		in.ProductId,
	).Scan(&productName, &cover, &price, &stock, &status)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.CartItemResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query product failed: %v", err)
		return nil, err
	}
	if status != 1 {
		return &order.CartItemResp{Code: errcode.BadRequest, Msg: "商品已下架"}, nil
	}

	var cartID int64
	var existingQuantity int32
	err = tx.QueryRowContext(l.ctx,
		"SELECT id, quantity FROM cart_item WHERE user_id = ? AND product_id = ? FOR UPDATE",
		userID, in.ProductId,
	).Scan(&cartID, &existingQuantity)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		if in.Quantity > stock {
			return &order.CartItemResp{Code: errcode.BadRequest, Msg: "库存不足"}, nil
		}
		result, err := tx.ExecContext(l.ctx,
			"INSERT INTO cart_item (user_id, product_id, quantity, selected) VALUES (?, ?, ?, 1)",
			userID, in.ProductId, in.Quantity,
		)
		if err != nil {
			l.Errorf("insert cart item failed: %v", err)
			return nil, err
		}
		cartID, err = result.LastInsertId()
		if err != nil {
			l.Errorf("read inserted cart item id failed: %v", err)
			return nil, err
		}
	case err != nil:
		l.Errorf("query cart item failed: %v", err)
		return nil, err
	default:
		quantity := existingQuantity + in.Quantity
		if quantity > maxQuantity {
			quantity = maxQuantity
		}
		if quantity > stock {
			return &order.CartItemResp{Code: errcode.BadRequest, Msg: "库存不足"}, nil
		}
		if _, err := tx.ExecContext(l.ctx,
			"UPDATE cart_item SET quantity = ?, selected = 1 WHERE id = ? AND user_id = ?",
			quantity, cartID, userID,
		); err != nil {
			l.Errorf("update cart item failed: %v", err)
			return nil, err
		}
	}

	info, err := findCartItem(l.ctx, tx, cartID, userID)
	if err != nil {
		l.Errorf("query cart item failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit cart transaction failed: %v", err)
		return nil, err
	}
	return &order.CartItemResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

type UpdateCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateCartLogic) UpdateCart(in *order.UpdateCartReq) (*order.CartItemResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.CartItemResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 || in.Quantity < 1 || in.Quantity > maxQuantity || (in.Selected != 0 && in.Selected != 1) {
		return &order.CartItemResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	tx, err := l.svcCtx.DB.BeginTx(l.ctx, nil)
	if err != nil {
		l.Errorf("begin cart transaction failed: %v", err)
		return nil, err
	}
	defer rollbackTx(l.Logger, tx)

	var (
		stock  int32
		status int32
	)
	err = tx.QueryRowContext(l.ctx,
		"SELECT COALESCE(p.stock, 0), COALESCE(p.status, 0) FROM cart_item ci LEFT JOIN looklook_product.product p ON p.id = ci.product_id AND p.deleted_at IS NULL WHERE ci.id = ? AND ci.user_id = ? FOR UPDATE",
		in.Id, userID,
	).Scan(&stock, &status)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &order.CartItemResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	case err != nil:
		l.Errorf("query cart item failed: %v", err)
		return nil, err
	}
	if status != 1 {
		return &order.CartItemResp{Code: errcode.BadRequest, Msg: "商品已下架"}, nil
	}
	if in.Quantity > stock {
		return &order.CartItemResp{Code: errcode.BadRequest, Msg: "库存不足"}, nil
	}

	if _, err := tx.ExecContext(l.ctx,
		"UPDATE cart_item SET quantity = ?, selected = ? WHERE id = ? AND user_id = ?",
		in.Quantity, in.Selected, in.Id, userID,
	); err != nil {
		l.Errorf("update cart item failed: %v", err)
		return nil, err
	}

	info, err := findCartItem(l.ctx, tx, in.Id, userID)
	if err != nil {
		l.Errorf("query cart item failed: %v", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		l.Errorf("commit cart transaction failed: %v", err)
		return nil, err
	}
	return &order.CartItemResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success), Data: info}, nil
}

type DeleteCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCartLogic {
	return &DeleteCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteCartLogic) DeleteCart(in *order.CartItemReq) (*order.BaseResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Id <= 0 {
		return &order.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	result, err := l.svcCtx.DB.ExecContext(l.ctx,
		"DELETE FROM cart_item WHERE id = ? AND user_id = ?", in.Id, userID,
	)
	if err != nil {
		l.Errorf("delete cart item failed: %v", err)
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		l.Errorf("read affected rows failed: %v", err)
		return nil, err
	}
	if affected == 0 {
		return &order.BaseResp{Code: errcode.NotFound, Msg: errcode.Msg(errcode.NotFound)}, nil
	}
	return &order.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}

type SelectCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectCartLogic {
	return &SelectCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SelectCartLogic) SelectCart(in *order.SelectCartReq) (*order.BaseResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.BaseResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}
	if in.Selected != 0 && in.Selected != 1 {
		return &order.BaseResp{Code: errcode.BadRequest, Msg: errcode.Msg(errcode.BadRequest)}, nil
	}

	if _, err := l.svcCtx.DB.ExecContext(l.ctx,
		"UPDATE cart_item SET selected = ? WHERE user_id = ?", in.Selected, userID,
	); err != nil {
		l.Errorf("select cart items failed: %v", err)
		return nil, err
	}
	return &order.BaseResp{Code: errcode.Success, Msg: errcode.Msg(errcode.Success)}, nil
}

type ListCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCartLogic {
	return &ListCartLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListCartLogic) ListCart(_ *order.EmptyReq) (*order.CartListResp, error) {
	userID, ok := auth.UserIDFromIncoming(l.ctx)
	if !ok {
		return &order.CartListResp{Code: errcode.Unauthorized, Msg: errcode.Msg(errcode.Unauthorized)}, nil
	}

	rows, err := l.svcCtx.DB.QueryContext(l.ctx,
		"SELECT "+cartSelectColumns+cartFromClause+" WHERE ci.user_id = ? ORDER BY ci.updated_at DESC",
		userID,
	)
	if err != nil {
		l.Errorf("query cart items failed: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			l.Errorf("close cart rows failed: %v", closeErr)
		}
	}()

	list := make([]*order.CartItemInfo, 0, 8)
	for rows.Next() {
		info, err := scanCartItem(rows)
		if err != nil {
			l.Errorf("scan cart item failed: %v", err)
			return nil, err
		}
		list = append(list, info)
	}
	if err := rows.Err(); err != nil {
		l.Errorf("iterate cart items failed: %v", err)
		return nil, err
	}

	return &order.CartListResp{
		Code:  errcode.Success,
		Msg:   errcode.Msg(errcode.Success),
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func rollbackTx(logger logx.Logger, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger.Errorf("rollback transaction failed: %v", err)
	}
}
