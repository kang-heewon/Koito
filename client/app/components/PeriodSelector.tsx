import { useEffect } from "react"

const PERIODS = ["day", "week", "month", "year", "all_time"]

interface Props {
    setter: Function
    current: string
    disableCache?: boolean
}

export default function PeriodSelector({ setter, current, disableCache = false }: Props) {
    const cacheScope = window.location.pathname

    const periodDisplay = (str: string) => {
        return str.split('_').map(w => w.split('').map((char, index) =>
            index === 0 ? char.toUpperCase() : char).join('')).join(' ')
    }

    const setPeriod = (val: string) => {
        setter(val)
        if (!disableCache) {
            localStorage.setItem('period_selection_' + cacheScope, val) 
        }  
    }

    useEffect(() => {
        if (!disableCache) {
            const cached = localStorage.getItem('period_selection_' + cacheScope);
            if (cached && PERIODS.includes(cached)) {
              setter(cached);
            }
        }
      }, [cacheScope, disableCache, setter]);

    return (
        <div className="flex gap-2 grow-0 text-sm sm:text-[16px]">
            <p>Showing stats for:</p>
            {PERIODS.map((p, i) => (
                <div key={`period_setter_${p}`}>
                    <button 
                        type="button"
                        className={`period-selector ${p === current ? 'color-fg' : 'color-fg-secondary'} ${i !== PERIODS.length - 1 ? 'pr-2' : ''}`}
                        onClick={() => setPeriod(p)}
                        disabled={p === current}
                    >
                        {periodDisplay(p)}
                    </button>
                    <span className="color-fg-secondary">
                        {i !== PERIODS.length - 1 ? '|' : ''}
                    </span>
                </div>
            ))}
        </div>
    )
}
