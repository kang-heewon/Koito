import ChartLayout from "./ChartLayout";
import { Link, useLoaderData, type LoaderFunctionArgs } from "react-router";
import { deleteListen, type Listen, type PaginatedResponse } from "api/api";
import { timeSince } from "~/utils/utils";
import ArtistLinks from "~/components/ArtistLinks";
import { useState } from "react";
import { useAppContext } from "~/providers/AppProvider";

export async function clientLoader({ request }: LoaderFunctionArgs) {
    const url = new URL(request.url);
    const page = url.searchParams.get("page") || "1";
    url.searchParams.set('page', page)

    const res = await fetch(
        `/apis/web/v1/listens?${url.searchParams.toString()}`
    );
    if (!res.ok) {
        throw new Response("Failed to load listens", { status: 500 });
    }

    const listens: PaginatedResponse<Listen> = await res.json();
    return { listens };
}

export default function Listens() {
    const { listens: initialData } = useLoaderData<{ listens: PaginatedResponse<Listen> }>();

    const [removedTimes, setRemovedTimes] = useState<Set<string>>(new Set())
    const [deleteError, setDeleteError] = useState('')
    const { user } = useAppContext()

    const buildRemovalKey = (page: number, listenTime: string) => `${page}:${listenTime}`

    const handleDelete = async (listen: Listen, page: number) => {
        setDeleteError('')

        try {
            const res = await deleteListen(listen)
            if (res.ok || (res.status >= 200 && res.status < 300)) {
                setRemovedTimes((prev) => {
                    const next = new Set(prev)
                    next.add(buildRemovalKey(page, listen.time))
                    return next
                })
            } else {
                const body = (await res.json().catch(() => null)) as { error?: string } | null
                setDeleteError(body?.error ?? `failed to delete listen (${res.status})`)
            }
        } catch (err) {
            setDeleteError(err instanceof Error ? err.message : "failed to delete listen")
        }
    }
  
  return (
        <ChartLayout
        title="Last Played"
        initialData={initialData}
        endpoint="listens"
        render={({ data, page, onNext, onPrev }) => (
            <div className="flex flex-col gap-5 text-sm md:text-[16px]">
            <div className="flex gap-15 mx-auto">
            <button type="button" className="default" onClick={onPrev} disabled={page <= 1}>
                Prev
            </button>
            <button type="button" className="default" onClick={onNext} disabled={!data.has_next_page}>
                Next
            </button>
            </div>
                <table className="-ml-4">
                    <tbody>
                        {data.items
                            .filter((item) => !removedTimes.has(buildRemovalKey(page, item.time)))
                            .map((item) => (
                            <tr key={`last_listen_${item.time}`} className="group hover:bg-[--color-bg-secondary]">
                                <td className="w-[17px] pr-2 align-middle">
                                    <button
                                        type="button"
                                        onClick={() => handleDelete(item, page)}
                                        className="opacity-0 group-hover:opacity-100 transition-opacity text-(--color-fg-tertiary) hover:text-(--color-error)"
                                        aria-label="Delete"
                                        hidden={user === null || user === undefined}
                                    >
                                        ×
                                    </button>
                                </td>
                                <td
                                    className="color-fg-tertiary pr-2 sm:pr-4 text-sm whitespace-nowrap w-0"
                                    title={new Date(item.time).toString()}
                                >
                                    {timeSince(new Date(item.time))}
                                </td>
                                <td className="text-ellipsis overflow-hidden max-w-[400px] sm:max-w-[600px]">
                                            <ArtistLinks artists={item.track.artists} /> –{' '}
                                    <Link
                                        className="hover:text-[--color-fg-secondary]"
                                        to={`/track/${item.track.id}`}
                                    >
                                        {item.track.title}
                                    </Link>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            <div className="flex gap-15 mx-auto">
                <button type="button" className="default" onClick={onPrev} disabled={page <= 1}>
                Prev
                </button>
                <button type="button" className="default" onClick={onNext} disabled={!data.has_next_page}>
                Next
                </button>
            </div>
            {deleteError !== '' && <p className="error text-center">{deleteError}</p>}
            </div>
        )}
        />
    );
}
