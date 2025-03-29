import { FC, useMemo } from "react"

interface IPill {
    title: string
    confidence: number
}

const Pill: FC<IPill> = ({ title, confidence }) => {
    const backgroundColor = useMemo(() => {
        if (confidence >= 80) return "bg-green-500"
        if (confidence >= 50) return "bg-yellow-500"
        return "bg-red-400"
    }, [confidence])

    return (
        <div className={`inline-flex items-center px-4 py-2 rounded-full text-white ${backgroundColor} transition`}>
            <span className="font-semibold">{title}</span>
            <span className="ml-2 text-sm opacity-90">({confidence}%)</span>
        </div>
    )
}

export default Pill