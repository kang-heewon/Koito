import { useState, type Dispatch, type SetStateAction } from "react";
import { Modal } from "./Modal";
import { replaceImage } from "api/api";
import { AsyncButton } from "../AsyncButton";

interface Props {
  type: string;
  id: number;
  musicbrainzId?: string;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export default function ImageReplaceModal({
  musicbrainzId,
  type,
  id,
  open,
  setOpen,
}: Props) {
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [suggestedImgLoading, setSuggestedImgLoading] = useState(true);

  const parseResponseError = async (r: Response): Promise<string> => {
    try {
      const body = (await r.json()) as { error?: string };
      if (body && typeof body.error === "string" && body.error.length > 0) {
        return body.error;
      }
    } catch {
      return `request failed (${r.status})`;
    }
    return `request failed (${r.status})`;
  };

  const doImageReplace = (url: string) => {
    setLoading(true);
    setError("");
    const formData = new FormData();
    formData.set(`${type.toLowerCase()}_id`, id.toString());
    formData.set("image_url", url);
    replaceImage(formData)
      .then(async (r) => {
        if (r.ok) {
          window.location.reload();
        } else {
          setError(await parseResponseError(r));
        }
      })
      .catch((err) =>
        setError(err instanceof Error ? err.message : "Failed to replace image")
      )
      .finally(() => setLoading(false));
  };

  const closeModal = () => {
    setOpen(false);
    setQuery("");
    setError("");
    setSuggestedImgLoading(true);
  };

  return (
    <Modal isOpen={open} onClose={closeModal}>
      <h2>Replace Image</h2>
      <div className="flex flex-col items-center">
        <input
          type="text"
          // i find my stupid a(n) logic to be a little silly so im leaving it in even if its not optimal
          placeholder={`Enter image URL, or drag-and-drop a local file`}
          className="w-full mx-auto fg bg rounded p-2"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        {query !== "" ? (
          <div className="flex gap-2 mt-4">
            <AsyncButton
              loading={loading}
              onClick={() => doImageReplace(query)}
            >
              Submit
            </AsyncButton>
          </div>
        ) : (
          ""
        )}
        {type === "Album" && musicbrainzId ? (
          <>
            <h3 className="mt-5">Suggested Image (Click to Apply)</h3>
            <button
              type="button"
              className="mt-4"
              disabled={loading}
              onClick={() =>
                doImageReplace(
                  `https://coverartarchive.org/release/${musicbrainzId}/front`
                )
              }
            >
              <div className={`relative`}>
                {suggestedImgLoading && (
                  <div className="absolute inset-0 flex items-center justify-center">
                    <div
                      className="animate-spin rounded-full border-2 border-gray-300 border-t-transparent"
                      style={{ width: 20, height: 20 }}
                    />
                  </div>
                )}
                <img
                  src={`https://coverartarchive.org/release/${musicbrainzId}/front`}
                  alt="Suggested album cover"
                  onLoad={() => setSuggestedImgLoading(false)}
                  onError={() => setSuggestedImgLoading(false)}
                  className={`block w-[130px] h-auto ${
                    suggestedImgLoading ? "opacity-0" : "opacity-100"
                  } transition-opacity duration-300`}
                />
              </div>
            </button>
          </>
        ) : (
          ""
        )}
        <p className="error">{error}</p>
      </div>
    </Modal>
  );
}
