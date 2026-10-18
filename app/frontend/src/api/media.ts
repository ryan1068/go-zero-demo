import { mediaRequest } from '@/utils/request'

export async function uploadImage(file: File): Promise<{ code: number; msg: string; url?: string }> {
  const form = new FormData()
  form.append('file', file)
  const res = await mediaRequest.post<{ code: number; msg: string; url?: string }>('/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data
}

export async function uploadImages(files: File[]): Promise<{ code: number; msg: string; urls: string[] }> {
  const form = new FormData()
  for (const file of files) {
    form.append('files', file)
  }
  const res = await mediaRequest.post<{ code: number; msg: string; urls: string[] }>('/uploads/batch', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data
}

export async function deleteMedia(name: string): Promise<{ code: number; msg: string }> {
  const res = await mediaRequest.delete<{ code: number; msg: string }>(
    `/admin/files/${encodeURIComponent(name)}`,
  )
  return res.data
}
