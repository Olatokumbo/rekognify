import { useCallback, useState } from 'react'
import { useDropzone } from 'react-dropzone'
import { generateS3PreSignedURL } from '../helpers'
import axios from 'axios'
import { useMainContext } from '../context/main-context'

const DropZone = () => {
    const [selectedFile, setSelectedFile] = useState<File | null>(null)
    const [preview, setPreview] = useState<string | null>(null)

    const { setUuid } = useMainContext()

    const onDrop = useCallback((acceptedFiles: File[]) => {
        const file = acceptedFiles[0]
        if (file) {
            setSelectedFile(file)
            setPreview(URL.createObjectURL(file))
        }
    }, [])

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: { 'image/*': [] },
        multiple: false,
    })

    const handleUpload = async () => {
        if (!selectedFile) {
            alert("Please select an image first!")
            return
        }

        try {
            const { id, url } = await generateS3PreSignedURL({
                filename: selectedFile.name,
                mimeType: selectedFile.type,
            })

            setUuid(id)

            console.log(url, id)

            const formData = new FormData()
            formData.append("file", selectedFile)

            const response = await axios.post(url, formData, {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
            })

            if (response.status === 200) {
                alert("Upload successful!")
            } else {
                alert("Upload failed.")
            }
        } catch (error) {
            console.error("Upload error:", error)
            alert("Something went wrong.")
        }
    }

    return (
        <div className="bg-white rounded-2xl p-3 m-3 w-[500px] min-h-[500px] flex flex-col items-center justify-center border-2 border-dashed border-gray-300 hover:border-orange-500 transition">
            <div {...getRootProps()} className="w-full h-full flex flex-col items-center justify-center cursor-pointer">
                <input {...getInputProps()} />
                {preview ? (
                    <img src={preview} className="w-full h-[400px] object-cover rounded-md" alt="Preview" />
                ) : (
                    <p className="text-gray-600 text-lg">
                        {isDragActive ? "Drop the image here..." : "Drag & drop an image here, or click to upload"}
                    </p>
                )}
            </div>
            <button
                onClick={handleUpload}
                className={`mt-4 px-6 py-2 rounded-md text-white transition shadow-lg 
                    ${selectedFile ? "bg-orange-500 hover:bg-orange-600" : "bg-gray-400 cursor-not-allowed"}`}
                disabled={!selectedFile}
            >
                Classify
            </button>
        </div>
    )
}

export default DropZone