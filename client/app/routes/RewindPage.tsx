import { ChevronLeft, ChevronRight } from "lucide-react";
import { average } from "color.js";
import { getRewindStats, imageUrl, type RewindStats } from "api/api";
import { useEffect, useState } from "react";
import { motion } from "motion/react";
import type { LoaderFunctionArgs } from "react-router";
import { useLoaderData, useLocation, useNavigate } from "react-router";
import Rewind from "~/components/rewind/Rewind";
import { getRewindParams } from "~/utils/utils";

const months = [
  "Full Year",
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

const fallbackAccentColor = "rgba(93, 211, 255, 0.3)";
const fallbackAccentGlow = "rgba(93, 211, 255, 0.16)";

export async function clientLoader({ request }: LoaderFunctionArgs) {
  const url = new URL(request.url);
  const fallbackParams = getRewindParams(url.searchParams);
  const parsedYear = Number.parseInt(url.searchParams.get("year") || "", 10);
  const parsedMonth = Number.parseInt(url.searchParams.get("month") || "", 10);
  const year = Number.isNaN(parsedYear) ? fallbackParams.year : parsedYear;
  const month = Number.isNaN(parsedMonth) ? fallbackParams.month : parsedMonth;

  const stats = await getRewindStats({ year, month });
  stats.title = `Your ${month === 0 ? "" : `${months[month]} `}${year} Rewind`;

  return { stats };
}

export function meta({ data }: { data?: { stats: RewindStats } }) {
  const pageTitle = `${data?.stats.title || "Rewind"} - Koito`;

  return [
    { title: pageTitle },
    { name: "description", content: pageTitle },
    { property: "og:title", content: pageTitle },
  ];
}

function NavigationControl({
  label,
  value,
  onPrev,
  onNext,
  prevDisabled,
  nextDisabled,
}: {
  label: string;
  value: string | number;
  onPrev: () => void;
  onNext: () => void;
  prevDisabled: boolean;
  nextDisabled: boolean;
}) {
  return (
    <div className="flex items-center gap-2 rounded-full border border-white/10 bg-[var(--color-bg)]/55 px-2 py-2 text-[var(--color-fg)] backdrop-blur-md sm:gap-3 sm:px-3">
      <span className="px-2 text-[0.65rem] font-semibold uppercase tracking-[0.22em] text-[var(--color-fg)]/55 sm:text-[0.7rem]">
        {label}
      </span>

      <button
        type="button"
        onClick={onPrev}
        className="flex h-9 w-9 items-center justify-center rounded-full border border-white/10 bg-white/5 transition-colors hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-35"
        disabled={prevDisabled}
        aria-label={`${label} 이전으로 이동`}
      >
        <ChevronLeft size={18} />
      </button>

      <span className="min-w-22 text-center text-sm font-semibold tracking-[0.08em] sm:min-w-28 sm:text-base">
        {value}
      </span>

      <button
        type="button"
        onClick={onNext}
        className="flex h-9 w-9 items-center justify-center rounded-full border border-white/10 bg-white/5 transition-colors hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-35"
        disabled={nextDisabled}
        aria-label={`${label} 다음으로 이동`}
      >
        <ChevronRight size={18} />
      </button>
    </div>
  );
}

export default function RewindPage() {
  const { stats } = useLoaderData() as { stats: RewindStats };
  const location = useLocation();
  const navigate = useNavigate();
  const [accentColor, setAccentColor] = useState(fallbackAccentColor);
  const [accentGlow, setAccentGlow] = useState(fallbackAccentGlow);

  const currentParams = new URLSearchParams(location.search);
  const { year, month } = getRewindParams(currentParams);
  const monthLabel = months[month];
  const now = new Date();

  useEffect(() => {
    const image = stats.top_artists[0]?.item?.image;
    if (!image) {
      setAccentColor(fallbackAccentColor);
      setAccentGlow(fallbackAccentGlow);
      return;
    }

    average(imageUrl(image, "small"), { amount: 1 })
      .then((color) => {
        const [red, green, blue] = color as unknown as number[];

        if (
          typeof red !== "number" ||
          typeof green !== "number" ||
          typeof blue !== "number"
        ) {
          setAccentColor(fallbackAccentColor);
          setAccentGlow(fallbackAccentGlow);
          return;
        }

        setAccentColor(`rgba(${red}, ${green}, ${blue}, 0.36)`);
        setAccentGlow(`rgba(${red}, ${green}, ${blue}, 0.18)`);
      })
      .catch(() => {
        setAccentColor(fallbackAccentColor);
        setAccentGlow(fallbackAccentGlow);
      });
  }, [stats]);

  const updateParams = (params: Record<string, string | null>) => {
    const nextParams = new URLSearchParams(location.search);

    for (const key in params) {
      const value = params[key];

      if (value === null) {
        nextParams.delete(key);
        continue;
      }

      nextParams.set(key, value);
    }

    navigate(`/rewind?${nextParams.toString()}`, { replace: false });
  };

  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const navigateMonth = (direction: "prev" | "next") => {
    let nextMonth = month;

    if (direction === "next") {
      nextMonth = month === 12 ? 0 : month + 1;
    } else {
      nextMonth = month === 0 ? 12 : month - 1;
    }

    scrollToTop();

    updateParams({
      year: String(year),
      month: String(nextMonth),
    });
  };

  const navigateYear = (direction: "prev" | "next") => {
    const nextYear = direction === "next" ? year + 1 : year - 1;

    scrollToTop();

    updateParams({
      year: String(nextYear),
      month: String(month),
    });
  };

  const prevMonthDisabled =
    new Date(year, month - 2) > now ||
    (now.getFullYear() === year && month === 1);
  const nextMonthDisabled = month >= now.getMonth() && year >= now.getFullYear();
  const prevYearDisabled = new Date(year - 1, month) > now;
  const nextYearDisabled =
    new Date(year + 1, month - 1) > now ||
    (month === 0 && now.getFullYear() === year + 1) ||
    (now.getMonth() === month - 1 && now.getFullYear() === year + 1);

  return (
    <main className="relative min-h-screen w-full text-[var(--color-fg)]">
      <div className="pointer-events-none fixed inset-x-0 top-0 z-30 px-4 pt-18 sm:px-8 sm:pt-8 lg:px-12">
        <div className="pointer-events-auto ml-auto flex w-full max-w-max flex-col gap-3 sm:flex-row sm:items-center">
          <NavigationControl
            label="Month"
            value={monthLabel}
            onPrev={() => navigateMonth("prev")}
            onNext={() => navigateMonth("next")}
            prevDisabled={prevMonthDisabled}
            nextDisabled={nextMonthDisabled}
          />

          <NavigationControl
            label="Year"
            value={year}
            onPrev={() => navigateYear("prev")}
            onNext={() => navigateYear("next")}
            prevDisabled={prevYearDisabled}
            nextDisabled={nextYearDisabled}
          />
        </div>
      </div>

      <motion.div
        key={`${year}-${month}`}
        initial={{ opacity: 0, y: 18 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.45, ease: [0.22, 1, 0.36, 1] }}
      >
        <Rewind
          stats={stats}
          monthLabel={monthLabel}
          year={year}
          accentColor={accentColor}
          accentGlow={accentGlow}
        />
      </motion.div>
    </main>
  );
}
