import { useQuery } from "@tanstack/react-query";
import { createApiKey, getApiKeys, type ApiKey } from "api/api";
import { useDeleteApiKey } from "../../hooks/useDeleteApiKey";
import { AsyncButton } from "../AsyncButton";
import { useEffect, useRef, useState } from "react";
import { Copy, Trash } from "lucide-react";

type CopiedState = {
    x: number;
    y: number;
    visible: boolean;
};

export default function ApiKeysModal() {
    const [input, setInput] = useState('')
    const [loading, setLoading ] = useState(false)
    const deleteApiKeyMutation = useDeleteApiKey()
    const [err, setError ] = useState<string>()
    const [displayData, setDisplayData] = useState<ApiKey[]>([])
    const [copied, setCopied] = useState<CopiedState | null>(null);
    const [expandedKey, setExpandedKey] = useState<string | null>(null);
    const textRefs = useRef<Record<string, HTMLButtonElement | null>>({});
    
    const handleRevealAndSelect = (key: string) => {
        setExpandedKey(key);
        setTimeout(() => {
            const el = textRefs.current[key];
            if (el) {
                const range = document.createRange();
                range.selectNodeContents(el);
                const sel = window.getSelection();
                sel?.removeAllRanges();
                sel?.addRange(range);
            }
        }, 0);
    };
        
    const { isPending, isError, data, error } = useQuery({ 
        queryKey: [
            'api-keys'
        ], 
        queryFn: () => {
            return getApiKeys();
        },
    });

    useEffect(() => {
        if (data) {
            setDisplayData(data)
        }
    }, [data])    

    if (isError) {
        return (
            <p className="error">Error: {error.message}</p>
        )
    }
    if (isPending) {
        return (
            <p>Loading...</p>
        )
    }

    const handleCopy = (e: React.MouseEvent<HTMLButtonElement>, text: string) => {
        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).catch(() => fallbackCopy(text));
        } else {
            fallbackCopy(text);
        }

        const parent = e.currentTarget.closest(".relative");
        if (!parent) {
            return;
        }

        const parentRect = parent.getBoundingClientRect();
        const buttonRect = e.currentTarget.getBoundingClientRect();
    
        setCopied({
            x: buttonRect.left - parentRect.left + buttonRect.width / 2,
            y: buttonRect.top - parentRect.top - 8,
            visible: true,
        });
    
        setTimeout(() => setCopied(null), 1500);
    };
    
    const fallbackCopy = (text: string) => {
        const textarea = document.createElement("textarea");
        textarea.value = text;
        textarea.style.position = "fixed"; // prevent scroll to bottom
        document.body.appendChild(textarea);
        textarea.focus();
        textarea.select();
        try {
            document.execCommand("copy");
        } catch {
            setError("Failed to copy API key");
        }
        document.body.removeChild(textarea);
    };

    const handleCreateApiKey = async () => {
        setError(undefined)
        if (input === "") {
            setError("a label must be provided")
            return
        }

        if (loading) {
            return
        }

        setLoading(true)

        try {
            const createdKey = await createApiKey(input)
            setDisplayData((prev) => [createdKey, ...prev])
            setInput('')
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to create API key")
        } finally {
            setLoading(false)
        }
    }

    const handleDeleteApiKey = async (id: number) => {
        setError(undefined)
        try {
            const r = await deleteApiKeyMutation.mutateAsync(id)
            if (r.ok) {
                setDisplayData((prev) => prev.filter((v) => v.id !== id))
            } else {
                const body = (await r.json().catch(() => null)) as { error?: string } | null
                setError(body?.error ?? "Failed to delete API key")
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to delete API key")
        }
    }

    return (
        <div className="">
        <h2>API Keys</h2>
        <div className="flex flex-col gap-4 relative">
            {displayData.map((v) => (
                <div key={v.key} className="flex gap-2"><button
                        type="button"
                        ref={el => {
                            textRefs.current[v.key] = el;
                        }}
                        onClick={() => handleRevealAndSelect(v.key)}
                        className={`bg p-3 rounded-md flex-grow text-left cursor-pointer ${
                            expandedKey === v.key ? '' : 'truncate'
                        }`}
                        style={{ whiteSpace: 'nowrap' }}
                        title={v.key}
                    >
                        {expandedKey === v.key ? v.key : `${v.key.slice(0, 8)}... ${v.label}`}
                    </button>
                    <button type="button" onClick={(e) => handleCopy(e, v.key)} className="large-button px-5 rounded-md" aria-label="Copy API key"><Copy size={16} /></button>
                    <AsyncButton loading={deleteApiKeyMutation.isPending} onClick={() => handleDeleteApiKey(v.id)} confirm><Trash size={16} /></AsyncButton>
                </div>
            ))}
            <div className="flex gap-2 w-3/5">
                <input
                    type="text"
                    placeholder="Add a label for a new API key"
                    className="mx-auto fg bg rounded-md p-3 flex-grow"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                />
                <AsyncButton loading={loading} onClick={handleCreateApiKey}>Create</AsyncButton>
            </div>
            {err && <p className="error">{err}</p>}
            {copied?.visible && (
                <div
                    style={{
                        position: "absolute",
                        top: copied.y,
                        left: copied.x,
                        transform: "translate(-50%, -100%)",
                    }}
                    className="pointer-events-none bg-black text-white text-sm px-2 py-1 rounded shadow-lg opacity-90 animate-fade"
                >
                    Copied!
                </div>
            )}
        </div>
        </div>
    )
}
