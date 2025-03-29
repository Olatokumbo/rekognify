/* eslint-disable @typescript-eslint/no-explicit-any */
import axios from "axios"

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

export const generateS3PreSignedURL = async (mimeType: string): Promise<PreSignedData> => {
    try {
      const response = await axios.post(`${api_url}/upload`, {
        mimeType
      })
      return response.data as PreSignedData
    } catch {
      throw new Error('Error Generating S3 Presigned URL')
    }
}


export const fetchImageInfo = async (filename: string, retries = 0): Promise<FileInfo> => {
  const MAX_RETRIES = 5;
  const RETRY_DELAY_MS = 1000; 
  
  try {
    const response = await axios.get(`${api_url}/info/${filename}`);
    return response.data as FileInfo;
  } catch (error: any) {
    if (error.response?.status === 409 && retries < MAX_RETRIES) {
      console.warn(`Retrying fetchImageInfo... Attempt ${retries + 1}`);
      await new Promise((resolve) => setTimeout(resolve, RETRY_DELAY_MS * 2 ** retries));
      return fetchImageInfo(filename, retries + 1);
    }
    throw error;
  }
};