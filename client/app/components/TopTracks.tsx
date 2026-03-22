import { useQuery, keepPreviousData } from "@tanstack/react-query";
import { getTopTracks, type GetItemsArgs, type TopRanked, type Track } from "api/api";
import { Link } from "react-router";
import TopListSkeleton from "./skeletons/TopListSkeleton";
import TopItemList from "./TopItemList";

interface Props {
  limit: number;
  period: string;
  artistId?: number;
  albumId?: number;
}

const TopTracks = (props: Props) => {
  const { isPending, isError, data, error, isFetching } = useQuery({
    queryKey: [
      "top-tracks",
      {
        limit: props.limit,
        period: props.period,
        artist_id: props.artistId,
        album_id: props.albumId,
        page: 0,
      },
    ],
    queryFn: ({ queryKey }) => getTopTracks(queryKey[1] as GetItemsArgs),
    placeholderData: keepPreviousData,
  });
  const rankedData = data as typeof data & { items: TopRanked<Track>[] };
  if (isPending || isFetching) {
    return (
      <div className="w-full sm:w-[300px]">
        <h2>Top Tracks</h2>
        <TopListSkeleton numItems={props.limit} />
      </div>
    );
  } else if (isError) {
    return (
      <div className="w-full sm:w-[300px]">
        <h2>Top Tracks</h2>
        <p className="error">Error: {error.message}</p>
      </div>
    );
  }
  if (!rankedData?.items) return null;

  let params = "";
  params += props.artistId ? `&artist_id=${props.artistId}` : "";
  params += props.albumId ? `&album_id=${props.albumId}` : "";

  return (
    <div className="w-full sm:w-[300px]">
      <h2 className="hover:underline">
        <Link to={`/chart/top-tracks?period=${props.period}${params}`}>
          Top Tracks
        </Link>
      </h2>
      <div className="max-w-full sm:max-w-[300px]">
        <TopItemList type="track" data={rankedData} />
        {rankedData.items.length < 1 ? "No tracks found for this period." : ""}
      </div>
    </div>
  );
};

export default TopTracks;
