import { replaceImage } from 'api/api';
import { useEffect } from 'react';

interface Props {
    itemType: string;
    onComplete: () => void;
}

export default function ImageDropHandler({ itemType, onComplete }: Props) {
  useEffect(() => {
    const parseResponseError = async (r: Response): Promise<string> => {
      try {
        const body = (await r.json()) as { error?: string };
        if (body && typeof body.error === 'string' && body.error.length > 0) {
          return body.error;
        }
      } catch {
        return `Upload failed (${r.status})`;
      }
      return `Upload failed (${r.status})`;
    };

    const handleDragOver = (e: DragEvent) => {
      e.preventDefault();
    };

    const handleDrop = async (e: DragEvent) => {
      e.preventDefault();
      if (!e.dataTransfer?.files.length) return;

      const imageFile = Array.from(e.dataTransfer.files).find((file) =>
        file.type.startsWith('image/')
      );
      if (!imageFile) return;

      const pathname = window.location.pathname;
      const segments = pathname.split('/').filter((segment) => segment !== '');
      const lastSegment = segments[segments.length - 1];
      const itemId = Number(lastSegment);
      if (!Number.isInteger(itemId) || itemId <= 0) {
        console.error('Upload failed: invalid route id');
        return;
      }

      const formData = new FormData();
      formData.append('image', imageFile);
      formData.append(`${itemType.toLowerCase()}_id`, String(itemId));

      try {
        const r = await replaceImage(formData);
        if (r.ok) {
          onComplete();
          return;
        }
        console.error(await parseResponseError(r));
      } catch (err) {
        console.error(
          `Upload failed: ${err instanceof Error ? err.message : String(err)}`
        );
      }
    };

    window.addEventListener('dragover', handleDragOver);
    window.addEventListener('drop', handleDrop);

    return () => {
      window.removeEventListener('dragover', handleDragOver);
      window.removeEventListener('drop', handleDrop);
    };
  }, [itemType, onComplete]);

  return null;
}
