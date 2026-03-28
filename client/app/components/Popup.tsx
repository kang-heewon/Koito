import React, { type PropsWithChildren, useEffect, useLayoutEffect, useState, useRef } from 'react';

interface Props {
    inner: React.ReactNode
    position: string
    space: number
    extraClasses?: string
    hint?: string
}

export default function Popup({ inner, position, space, extraClasses, children }: PropsWithChildren<Props>) {
    const [isVisible, setIsVisible] = useState(false);
    const [showPopup, setShowPopup] = useState(true);
    const tooltipRef = useRef<HTMLDivElement>(null);
    const [adjustedLeft, setAdjustedLeft] = useState(70 + space);

    useEffect(() => {
        const mediaQuery = window.matchMedia('(min-width: 640px)');

        const handleChange = (e: MediaQueryListEvent) => {
            setShowPopup(e.matches);
        };

        setShowPopup(mediaQuery.matches);

        mediaQuery.addEventListener('change', handleChange);
        return () => mediaQuery.removeEventListener('change', handleChange);
    }, []);

    useLayoutEffect(() => {
        if (tooltipRef.current && showPopup) {
            const tooltipWidth = tooltipRef.current.offsetWidth;
            const baseLeft = 70 + space;

            const tooltipLeftEdge = baseLeft - tooltipWidth / 2;

            if (tooltipLeftEdge < 0) {
                setAdjustedLeft(tooltipWidth / 2);
            } else {
                setAdjustedLeft(baseLeft);
            }
        }
    }, [showPopup, space]);

    let positionClasses = '';
    let spaceCSS: React.CSSProperties = {};
    if (position === 'top') {
        positionClasses = `top-${space} -bottom-2 -translate-y-1/2 -translate-x-1/2`;
    } else if (position === 'right') {
        positionClasses = `bottom-1 -translate-x-1/2`;
        spaceCSS = { left: adjustedLeft };
    }

    return (
        <div
            className="relative"
            onMouseEnter={() => setIsVisible(true)}
            onMouseLeave={() => setIsVisible(false)}
        >
            {children}
            {showPopup && (
                <div
                    ref={tooltipRef}
                    className={`
                    absolute
                    ${positionClasses}
                    ${extraClasses ?? ''}
                    bg-(--color-bg) color-fg border-1 border-(--color-bg-tertiary)
                    px-3 py-2 rounded-lg
                    transition-opacity duration-100
                    ${isVisible ? 'opacity-100' : 'opacity-0 pointer-events-none'}
                    z-50 text-center
                    flex
                `}
                    style={spaceCSS}
                >
                    {inner}
                </div>
            )}
        </div>
    );
}
