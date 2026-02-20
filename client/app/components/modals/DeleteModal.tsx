import { deleteItem } from "api/api"
import { AsyncButton } from "../AsyncButton"
import { Modal } from "./Modal"
import { useNavigate } from "react-router"
import { useState } from "react"

interface Props {
    open: boolean
    setOpen: Function
    title: string,
    id: number,
    type: string
}

export default function DeleteModal({ open, setOpen, title, id, type }: Props) {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')
    const navigate = useNavigate()

    const doDelete = async () => {
        if (loading) {
            return
        }

        setError('')
        setLoading(true)

        try {
            const r = await deleteItem(type.toLowerCase(), id)
            if (r.ok) {
                navigate('/')
            } else {
                const body = (await r.json().catch(() => null)) as { error?: string } | null
                setError(body?.error ?? `failed to delete ${type.toLowerCase()}`)
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : `failed to delete ${type.toLowerCase()}`)
        } finally {
            setLoading(false)
        }
    }

    return (
        <Modal isOpen={open} onClose={() => setOpen(false)}>
            <h2>Delete "{title}"?</h2>
            <p>This action is irreversible!</p>
            <div className="flex flex-col mt-3 items-center">
                <AsyncButton loading={loading} onClick={doDelete}>Yes, Delete It</AsyncButton>
                {error !== '' && <p className="error mt-3">{error}</p>}
            </div>
        </Modal>
    )
}
