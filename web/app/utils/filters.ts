import { intervalToDuration, formatDuration } from 'date-fns'
import { parsePhoneNumber, isValidPhoneNumber } from 'libphonenumber-js'
import type { Locale } from 'date-fns'

export function formatPhoneNumber(value: string): string {
  if (!value || typeof value !== 'string') {
    return value ?? ''
  }
  try {
    if (!isValidPhoneNumber(value)) {
      return value
    }
    const phoneNumber = parsePhoneNumber(value)
    if (phoneNumber) {
      return phoneNumber.formatInternational()
    }
  } catch {
    return value
  }
  return value
}

export function phoneCountry(value: string): string {
  if (!value || typeof value !== 'string') return 'Earth'
  try {
    const phoneNumber = parsePhoneNumber(value)
    if (phoneNumber && phoneNumber.country) {
      const regionNames = new Intl.DisplayNames(undefined, { type: 'region' })
      return regionNames.of(phoneNumber.country) ?? 'Earth'
    }
  } catch {
    return 'Earth'
  }
  return 'Earth'
}

export function formatTimestamp(value: string): string {
  if (!value) return ''
  try {
    const d = new Date(value)
    if (isNaN(d.getTime())) return ''
    return d.toLocaleString()
  } catch {
    return ''
  }
}

export function formatMoney(value: string | number): string {
  try {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL',
    }).format(typeof value === 'string' ? parseFloat(value) : value)
  } catch {
    return String(value ?? '')
  }
}

export function formatDecimal(value: string | number): string {
  try {
    return new Intl.NumberFormat('pt-BR', {
      style: 'decimal',
    }).format(typeof value === 'string' ? parseFloat(value) : value)
  } catch {
    return String(value ?? '')
  }
}

export function formatBillingPeriod(value: string): string {
  if (!value) return ''
  try {
    const d = new Date(value)
    if (isNaN(d.getTime())) return ''
    return d.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
    })
  } catch {
    return ''
  }
}

export function formatBillingPeriodDateOrdinal(value: string): string {
  if (!value) return ''
  try {
    const date = new Date(value)
    if (isNaN(date.getTime())) return ''
    const day = date.getDate()
    const month = date.toLocaleDateString('en-US', { month: 'long' })
    const year = date.getFullYear()

    const suffix =
      day % 10 === 1 && day !== 11
        ? 'st'
        : day % 10 === 2 && day !== 12
          ? 'nd'
          : day % 10 === 3 && day !== 13
            ? 'rd'
            : 'th'

    return `${month} ${day}<sup>${suffix}</sup> ${year}`
  } catch {
    return ''
  }
}

export interface BillingPeriodDateOrdinalParts {
  leading: string
  suffix: string
  trailing: string
}

export function formatBillingPeriodDateOrdinalParts(
  value: string,
): BillingPeriodDateOrdinalParts {
  if (!value) return { leading: '', suffix: '', trailing: '' }
  try {
    const date = new Date(value)
    if (isNaN(date.getTime())) return { leading: '', suffix: '', trailing: '' }
    const day = date.getDate()
    const month = date.toLocaleDateString('en-US', { month: 'long' })
    const year = date.getFullYear()

    const suffix =
      day % 10 === 1 && day !== 11
        ? 'st'
        : day % 10 === 2 && day !== 12
          ? 'nd'
          : day % 10 === 3 && day !== 13
            ? 'rd'
            : 'th'

    return { leading: `${month} ${day}`, suffix, trailing: ` ${year}` }
  } catch {
    return { leading: '', suffix: '', trailing: '' }
  }
}

export function humanizeTime(value: string, locale?: Locale): string {
  if (!value) return ''
  try {
    const start = new Date(value)
    if (isNaN(start.getTime())) return ''
    const durations = intervalToDuration({
      start,
      end: new Date(),
    })
    return formatDuration(durations, { locale })
  } catch {
    return ''
  }
}

export function startsWithLetter(value: string): boolean {
  return /^[a-zA-Z]/.test(value)
}

export function humanizeTimeShort(date: string, locale?: Locale) {
  try {
    return formatDistanceToNow(new Date(date), { addSuffix: true, locale })
  } catch {
    return ''
  }
}
