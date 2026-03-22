import { ChevronLeft, ChevronRight } from "lucide-react";
import { average } from "color.js";
import { getRewindStats, imageUrl, type RewindStats } from "api/api";
import { useEffect, useState } from "react";
import { motion } from "motion/react";
import type { LoaderFunctionArgs } from "react-router";
import { redirect, useLoaderData, useNavigate } from "react-router";
import Rewind from "~/components/rewind/Rewind";

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

const getRewindPath = (year: number, month: number) => {
  if (month === 0) {
    return `/rewind/${year}`;
  }

  return `/rewind/${year}/${month}`;
};

const getNumericValue = (value?: string | null) => {
  if (!value) {
    return null;
  }

  const parsedValue = Number.parseInt(value, 10);

  if (Number.isNaN(parsedValue)) {
    return null;
  }

  return parsedValue;
};

const getLoaderRewindParams = ({ params, request }: LoaderFunctionArgs) => {
  const url = new URL(request.url);
  const searchYear = getNumericValue(url.searchParams.get("year"));
  const searchMonth = getNumericValue(url.searchParams.get("month"));
  const hasLegacySearchParams =
    url.searchParams.has("year") || url.searchParams.has("month");

  if (hasLegacySearchParams) {
    const targetYear = searchYear ?? new Date().getFullYear();
    const targetMonth = searchMonth ?? 0;

    return {
      month: targetMonth,
      redirectTo: getRewindPath(targetYear, targetMonth),
      year: targetYear,
    };
  }

  const currentYear = new Date().getFullYear();
  const year = getNumericValue(params.year) ?? currentYear;
  const month = getNumericValue(params.month) ?? 0;

  return {
    month,
    redirectTo: null,
    year,
  };
};

export async function clientLoader(args: LoaderFunctionArgs) {
  const { year, month, redirectTo } = getLoaderRewindParams(args);

  if (redirectTo) {
    throw redirect(redirectTo);
  }

  const stats = await getRewindStats({ year, month });
  stats.title = `Your ${month === 0 ? "" : `${months[month]} `}${year} Rewind`;

  return { month, stats, year };
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
  const { month, stats, year } = useLoaderData() as {
    month: number;
    stats: RewindStats;
    year: number;
  };
  const navigate = useNavigate();
  const [accentColor, setAccentColor] = useState(fallbackAccentColor);
  const [accentGlow, setAccentGlow] = useState(fallbackAccentGlow);

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

  const navigateToRewind = (nextYear: number, nextMonth: number) => {
    navigate(getRewindPath(nextYear, nextMonth), { replace: false });
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

    navigateToRewind(year, nextMonth);
  };

  const navigateYear = (direction: "prev" | "next") => {
    const nextYear = direction === "next" ? year + 1 : year - 1;

    scrollToTop();

    navigateToRewind(nextYear, month);
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
