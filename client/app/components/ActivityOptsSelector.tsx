import { ChevronDown, ChevronUp } from "lucide-react";
import { useEffect, useState } from "react";

const STEP_PERIODS = ["day", "week", "month"];
const RANGE_PERIODS = [105, 182, 364];

interface Props {
    stepSetter: (value: string) => void;
    currentStep: string;
    rangeSetter: (value: number) => void;
    currentRange: number;
    disableCache?: boolean;
}

export default function ActivityOptsSelector({
    stepSetter,
    currentStep,
    rangeSetter,
    currentRange,
    disableCache = false,
}: Props) {
    const [collapsed, setCollapsed] = useState(true);
    const cacheScope = window.location.pathname;

    const setMenuOpen = (val: boolean) => {
        setCollapsed(val)
        if (!disableCache) {
            localStorage.setItem('activity_configuring_' + cacheScope, String(!val));
        }
    }

    const setStep = (val: string) => {
        stepSetter(val);
        if (!disableCache) {
            localStorage.setItem('activity_step_' + cacheScope, val);
        }
    };

    const setRange = (val: number) => {
        rangeSetter(val);
        if (!disableCache) {
            localStorage.setItem('activity_range_' + cacheScope, String(val));
        }
    };

    useEffect(() => {
        if (!disableCache) {
            const cachedRange = localStorage.getItem('activity_range_' + cacheScope);
            if (cachedRange !== null) {
                const parsed = parseInt(cachedRange, 10);
                if (!Number.isNaN(parsed)) {
                    rangeSetter(parsed);
                }
            }

            const cachedStep = localStorage.getItem('activity_step_' + cacheScope);
            if (cachedStep && STEP_PERIODS.includes(cachedStep)) {
                stepSetter(cachedStep);
            }

            const cachedConfiguring = localStorage.getItem('activity_configuring_' + cacheScope);
            if (cachedConfiguring !== null) {
                setCollapsed(cachedConfiguring !== "true");
            }
        }
    }, [cacheScope, disableCache, rangeSetter, stepSetter]);

    return (
        <div className="relative w-full">
            <button
                type="button"
                onClick={() => setMenuOpen(!collapsed)}
                className="absolute left-[75px] -top-9 text-muted hover:color-fg transition"
                title="Toggle options"
                aria-expanded={!collapsed}
            >
                {collapsed ? <ChevronDown size={18} /> : <ChevronUp size={18} />}
            </button>

            <div
                className={`overflow-hidden transition-[max-height,opacity] duration-250 ease ${
                    collapsed ? 'max-h-0 opacity-0' : 'max-h-[100px] opacity-100'
                }`}
            >
                <div className="flex flex-wrap gap-4 mt-1 text-sm">
                    <div className="flex items-center gap-1">
                        <span className="text-muted">Step:</span>
                        {STEP_PERIODS.map((p) => (
                            <button
                                type="button"
                                key={p}
                                className={`px-1 rounded transition ${
                                    p === currentStep ? 'color-fg font-medium' : 'color-fg-secondary hover:color-fg'
                                }`}
                                onClick={() => setStep(p)}
                                disabled={p === currentStep}
                                aria-pressed={p === currentStep}
                            >
                                {p}
                            </button>
                        ))}
                    </div>

                    <div className="flex items-center gap-1">
                        <span className="text-muted">Range:</span>
                        {RANGE_PERIODS.map((r) => (
                            <button
                                type="button"
                                key={r}
                                className={`px-1 rounded transition ${
                                    r === currentRange ? 'color-fg font-medium' : 'color-fg-secondary hover:color-fg'
                                }`}
                                onClick={() => setRange(r)}
                                disabled={r === currentRange}
                                aria-pressed={r === currentRange}
                            >
                                {r}
                            </button>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
