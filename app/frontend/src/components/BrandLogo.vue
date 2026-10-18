<script setup lang="ts">
withDefaults(
  defineProps<{
    subtitle?: string
    size?: 'md' | 'lg'
    light?: boolean
  }>(),
  { subtitle: '', size: 'md', light: false },
)
</script>

<template>
  <div class="brand-logo" :class="[size, { light }]">
    <div class="seal" aria-hidden="true">
      <span>拾</span>
      <span>光</span>
    </div>
    <div class="brand-copy">
      <strong>拾光严选</strong>
      <small v-if="subtitle">{{ subtitle }}</small>
    </div>
  </div>
</template>

<style scoped>
.brand-logo {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  line-height: 1;
}

.seal {
  width: 42px;
  height: 42px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1px;
  border: 1px solid rgba(247, 239, 240, 0.65);
  border-radius: 7px;
  background: var(--brand);
  color: #f7eff0;
  font-family: var(--font-serif);
  font-size: 15px;
  letter-spacing: 1px;
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.12);
}

.brand-copy {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.brand-copy strong {
  font-family: var(--font-serif);
  font-size: 21px;
  letter-spacing: 3px;
  color: var(--ink);
  font-weight: 700;
}

.brand-copy small {
  color: var(--muted);
  font-size: 11px;
  letter-spacing: 2px;
}

.brand-logo.lg .seal {
  width: 58px;
  height: 58px;
  font-size: 20px;
  border-radius: 9px;
}

.brand-logo.lg .brand-copy strong {
  font-size: 30px;
  letter-spacing: 5px;
}

.brand-logo.lg .brand-copy small {
  font-size: 13px;
}

.brand-logo.light .brand-copy strong,
.brand-logo.light .brand-copy small {
  color: #f7eff0;
}

.brand-logo.light .brand-copy small {
  color: rgba(247, 239, 240, 0.72);
}
</style>
