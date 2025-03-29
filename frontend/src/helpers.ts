import axios from "axios"

interface UploadParams {
    filename: string
    mimeType: string
}

interface PreSignedData {
    url: string,
    id: string
}


interface Label {
    category: string;
    confidence: number;
    name: string;
  }
  
  interface FileInfo {
    filename: string;
    url: string;
    labels: Label[];
  }

const api_url = import.meta.env.VITE_API_URL

export const generateS3PreSignedURL = async ({ filename, mimeType }: UploadParams): Promise<PreSignedData> => {
    try {
      const response = await axios.post(`${api_url}/upload`, {
        filename,
        mimeType
      })
      return response.data as PreSignedData
    } catch {
      throw new Error('Error Generating S3 Presigned URL')
    }
}

export const fetchImageInfo = async (filename: string) =>{
    const response = await axios.get(`${api_url}/info/${filename}`)
    return response.data as FileInfo
}