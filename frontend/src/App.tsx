import Navbar from './components/navbar'
import Footer from './components/footer'
import DropZone from './components/dropzone'
import Pill from './components/pill'
import { useEffect, useState } from 'react'
import { fetchImageInfo } from './helpers'
import { useMainContext } from './context/main-context'

function App() {
  const { uuid } = useMainContext()
  const [labels, setLabels] = useState<{ title: string, confidence: number }[] | null>(null)

  useEffect(() => {
    (async () => {
      if (uuid) {
        const { labels } = await fetchImageInfo(uuid)
        setLabels(labels.map((label) => {
          return {
            title: label.name,
            confidence: label.confidence
          }
        }))
      }
    })()
  }, [uuid]);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <div className="flex-1 flex p-6 justify-center items-center">
        <div className='flex justify-center items-center bg-[#1E1E1E] min-w-5xl rounded-2xl p-6 flex-col border-2 transition shadow-2xl lg:flex-row'>
          <DropZone />
          <div className='flex flex-wrap justify-center gap-4 p-4 max-w-[400px] overflow-auto'>
            {labels && labels.map((label, index) => (
              <Pill key={index} title={label.title} confidence={label.confidence} />
            ))}
          </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default App