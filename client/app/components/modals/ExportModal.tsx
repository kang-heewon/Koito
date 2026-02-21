import { useState } from "react";
import { AsyncButton } from "../AsyncButton";
import { getExport } from "api/api";

export default function ExportModal() {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')

    const handleExport = async () => {
        if (loading) {
            return
        }

        setError('')
        setLoading(true)

        try {
            const blob = await getExport()
            const url = window.URL.createObjectURL(blob)
            const a = document.createElement("a")
            a.href = url
            a.download = "koito_export.json"
            document.body.appendChild(a)
            a.click()
            a.remove()
            window.URL.revokeObjectURL(url)
        } catch (err) {
            setError(err instanceof Error ? err.message : "failed to export data")
        } finally {
            setLoading(false)
        }
    }

    return (
        <div>
            <h2>Export</h2>
            <AsyncButton loading={loading} onClick={() => { void handleExport() }}>Export Data</AsyncButton>
            {error && <p className="error">{error}</p>}
        </div>
    )
}
