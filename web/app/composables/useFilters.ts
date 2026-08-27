import {
  formatPhoneNumber,
  phoneCountry,
  formatTimestamp,
  formatMoney,
  formatDecimal,
  formatBillingPeriod,
  formatBillingPeriodDateOrdinal,
  humanizeTime,
  humanizeTimeShort,
  startsWithLetter,
} from '../utils/filters'
import { capitalize } from '../utils/capitalize'

import { ptBR, enUS } from 'date-fns/locale'

export function useFilters() {
  const { locale } = useI18n()

  const getDateFnsLocale = () => {
    return locale.value === 'pt-BR' ? ptBR : enUS
  }

  return {
    formatPhoneNumber,
    phoneCountry,
    formatTimestamp,
    formatMoney,
    formatDecimal,
    formatBillingPeriod,
    humanizeTimeShort: (date: string) =>
      humanizeTimeShort(date, getDateFnsLocale()),
    formatBillingPeriodDateOrdinal,
    humanizeTime: (value: string) => humanizeTime(value, getDateFnsLocale()),
    startsWithLetter,
    capitalize,
  }
}
