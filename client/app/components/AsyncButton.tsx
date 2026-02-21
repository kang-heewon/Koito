import React, { useEffect, useRef, useState } from "react"

type Props = {
    children: React.ReactNode
    onClick: () => void
    loading?: boolean
    disabled?: boolean
    confirm?: boolean
}

export function AsyncButton(props: Props) {
    const [awaitingConfirm, setAwaitingConfirm] = useState(false)
    const confirmTimeoutRef = useRef<number | null>(null)

    useEffect(() => {
        if (props.loading || props.disabled) {
            setAwaitingConfirm(false)
        }
    }, [props.loading, props.disabled])

    useEffect(() => {
        return () => {
            if (confirmTimeoutRef.current !== null) {
                window.clearTimeout(confirmTimeoutRef.current)
            }
        }
    }, [])

    const handleClick = () => {
        if (props.confirm) {
            if (!awaitingConfirm) {
                setAwaitingConfirm(true)
                if (confirmTimeoutRef.current !== null) {
                    window.clearTimeout(confirmTimeoutRef.current)
                }
                confirmTimeoutRef.current = window.setTimeout(() => {
                    setAwaitingConfirm(false)
                    confirmTimeoutRef.current = null
                }, 3000)
                return
            }

            if (confirmTimeoutRef.current !== null) {
                window.clearTimeout(confirmTimeoutRef.current)
                confirmTimeoutRef.current = null
            }
            setAwaitingConfirm(false)
        }

        props.onClick()
    }

    return (
        <button
            type="button"
            onClick={handleClick}
            disabled={props.loading || props.disabled}
            className={`relative px-5 py-2 rounded-md large-button flex disabled:opacity-50 items-center`}
        >
            <span className={props.loading ? 'invisible' : 'visible'}>
                {awaitingConfirm ? 'Are you sure?' : props.children}
            </span>
            {props.loading && (
                <span className="absolute inset-0 flex items-center justify-center">
                    <span className="animate-spin h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
                </span>
            )}
        </button>
    )
}
