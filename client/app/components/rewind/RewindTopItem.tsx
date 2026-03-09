type Ranked<T> = {
  item: T;
  rank: number;
  listen_count: number;
  time_listened: number;
};

type TopItemProps<T> = {
  title: string;
  imageSrc: string;
  items: Ranked<T>[];
  getLabel: (item: T) => string;
  includeTime?: boolean;
};

export default function RewindTopItem<T extends { id: string | number }>({
  title,
  imageSrc,
  items,
  getLabel,
  includeTime,
}: TopItemProps<T>) {
  const [top, ...rest] = items;

  if (!top) {
    return null;
  }

  return (
    <div className="flex flex-col gap-5 sm:flex-row">
      <div className="rewind-top-item-image">
        <img className="max-h-48 max-w-48" src={imageSrc} alt={title} />
      </div>

      <div className="flex flex-col gap-1">
        <h4 className="-mb-1">{title}</h4>

        <div className="flex items-center gap-2">
          <div className="mb-2 flex flex-col items-start">
            <h2>{getLabel(top.item)}</h2>
            <span className="-mt-3 text-sm text-(--color-fg-tertiary)">
              {`${top.listen_count} plays`}
              {includeTime ? ` (${Math.floor(top.time_listened / 60)} minutes)` : ""}
            </span>
          </div>
        </div>

        {rest.map((entry) => (
          <div key={entry.item.id} className="text-sm">
            {getLabel(entry.item)}
            <span className="text-(--color-fg-tertiary)">
              {` - ${entry.listen_count} plays`}
              {includeTime
                ? ` (${Math.floor(entry.time_listened / 60)} minutes)`
                : ""}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
