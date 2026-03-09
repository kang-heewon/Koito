interface Props {
  figure: string;
  text: string;
}

export default function RewindStatText(props: Props) {
  return (
    <div className="flex items-baseline gap-1.5">
      <div className="w-23 shrink-0 text-end">
        <span
          className="
            relative inline-block
            text-2xl font-semibold
          "
        >
          <span
            className="
              absolute inset-0
              -translate-x-2 translate-y-8
              z-0 h-1
              bg-(--color-primary)
            "
            aria-hidden
          />
          <span className="relative z-1">{props.figure}</span>
        </span>
      </div>
      <span className="text-sm">{props.text}</span>
    </div>
  );
}
