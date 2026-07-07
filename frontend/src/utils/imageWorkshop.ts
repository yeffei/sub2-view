import type { CustomMenuItem } from '@/types'

export const IMAGE_WORKSHOP_MENU_ID = 'image-workshop'

export function findImageWorkshopMenuItem(items: CustomMenuItem[] | undefined | null): CustomMenuItem | null {
  return (items ?? []).find((item) => (
    item.id === IMAGE_WORKSHOP_MENU_ID
    && item.visibility === 'user'
    && (Boolean(item.url?.trim()) || Boolean(item.page_slug?.trim()))
  )) ?? null
}

export function isImageWorkshopMenuItem(item: CustomMenuItem): boolean {
  return item.id === IMAGE_WORKSHOP_MENU_ID
}
