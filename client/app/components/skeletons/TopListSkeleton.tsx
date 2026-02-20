interface Props {
    numItems: number
}

export default function TopListSkeleton({ numItems }: Props) {
    const skeletonItems = Array.from({ length: numItems }, (_, idx) => `top-list-skeleton-${idx + 1}`)

    return (
        <div className="w-[300px] animate-pulse" aria-hidden="true">
            {skeletonItems.map((itemKey) => (
                <div key={itemKey} className="flex items-center gap-2 mb-[4px]">
                    <div className="w-[40px] h-[40px] bg-(--color-bg-tertiary) rounded"></div>
                    <div>
                        <div className="h-[14px] w-[150px] bg-(--color-bg-tertiary) rounded"></div>
                        <div className="h-[12px] w-[60px] bg-(--color-bg-tertiary) rounded"></div>
                    </div>
                </div>
            ))}
        </div>
    )
}
