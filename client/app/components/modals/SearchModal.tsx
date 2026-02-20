import { useEffect, useState, type Dispatch, type SetStateAction } from "react";
import { Modal } from "./Modal";
import { search, type SearchResponse } from "api/api";
import SearchResults from "../SearchResults";

interface Props {
    open: boolean 
    setOpen: Dispatch<SetStateAction<boolean>>
}

export default function SearchModal({ open, setOpen }: Props) {
    const [query, setQuery] = useState('');
    const [data, setData] = useState<SearchResponse>();
    const [debouncedQuery, setDebouncedQuery] = useState(query);
    const [error, setError] = useState('');

    const closeSearchModal = () => {
        setOpen(false)
        setQuery('')
        setData(undefined)
        setError('')
    }

    useEffect(() => {
        const handler = setTimeout(() => {
            setDebouncedQuery(query);
            if (query === '') {
                setData(undefined)
                setError('')
            }
        }, 300);

        return () => {
            clearTimeout(handler);
        };
    }, [query]);

    useEffect(() => {
        let active = true;

        if (debouncedQuery) {
            search(debouncedQuery).then((r) => {
                if (!active) {
                    return
                }
                setError('')
                setData(r);
            }).catch((err) => {
                if (!active) {
                    return
                }
                setData(undefined)
                setError(err instanceof Error ? err.message : 'Search failed')
            });
        }

        return () => {
            active = false
        }
    }, [debouncedQuery]);

    return (
        <Modal isOpen={open} onClose={closeSearchModal}>
            <h2>Search</h2>
            <div className="flex flex-col items-center">
                <input
                    type="text"
                    placeholder="Search for an artist, album, or track"
                    className="w-full mx-auto fg bg rounded p-2"
                    onChange={(e) => setQuery(e.target.value)}
                />
                <div className="h-3/4 w-full">
                <SearchResults data={data} onSelect={closeSearchModal}/>
                </div>
                {error ? <p className="error mt-3">{error}</p> : null}
            </div>
        </Modal>
    )
}
