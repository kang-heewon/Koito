import AnimatedCounter from "../recap/AnimatedCounter";

interface Props {
  title: string;
  value: number;
  description: string;
  prefix?: string;
  suffix?: string;
}

export default function RewindStatText({
  title,
  value,
  description,
  prefix,
  suffix,
}: Props) {
  return (
    <div className="flex min-h-72 flex-col justify-between rounded-[32px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/80 px-6 py-7 text-[var(--color-fg)] backdrop-blur-sm sm:px-8 sm:py-9">
      <div>
        <div className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--color-primary)]/78">
          {title}
        </div>
        <div className="header-font mt-6 text-5xl font-semibold leading-none tracking-[-0.05em] sm:text-6xl">
          <AnimatedCounter value={value} prefix={prefix} suffix={suffix} />
        </div>
      </div>

      <p className="mt-8 text-base leading-7 text-[var(--color-fg)]/72">{description}</p>
    </div>
  );
}
