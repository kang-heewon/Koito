import { type Ranked } from "api/api";
import { TopItemCard } from "../recap";

type TopItemProps<T extends { id: string | number }> = {
  eyebrow: string;
  title: string;
  description: string;
  items: Ranked<T>[];
  getName: (item: T) => string;
  getImage: (item: T) => string | undefined;
  getSubtitle: (entry: Ranked<T>) => string | undefined;
  emptyState: string;
};

export default function RewindTopItem<T extends { id: string | number }>({
  eyebrow,
  title,
  description,
  items,
  getName,
  getImage,
  getSubtitle,
  emptyState,
}: TopItemProps<T>) {
  return (
    <div className="space-y-8">
      <div className="max-w-3xl">
        <div className="text-xs font-semibold uppercase tracking-[0.28em] text-[var(--color-primary)]/78">
          {eyebrow}
        </div>
        <h2 className="header-font mt-4 text-4xl font-semibold tracking-[-0.04em] text-[var(--color-fg)] sm:text-5xl lg:text-6xl">
          {title}
        </h2>
        <p className="mt-4 text-base leading-7 text-[var(--color-fg)]/72 sm:text-lg">
          {description}
        </p>
      </div>

      {items.length === 0 ? (
        <div className="flex min-h-64 items-center justify-center rounded-[32px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/80 px-6 py-8 text-center text-base text-[var(--color-fg)]/65 backdrop-blur-sm sm:px-8">
          {emptyState}
        </div>
      ) : (
        <div className="grid gap-4 xl:grid-cols-2">
          {items.slice(0, 5).map((entry) => (
            <TopItemCard
              key={entry.item.id}
              rank={entry.rank}
              name={getName(entry.item)}
              imageUrl={getImage(entry.item)}
              subtitle={getSubtitle(entry)}
              plays={entry.listen_count}
            />
          ))}
        </div>
      )}
    </div>
  );
}
