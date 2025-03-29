import captureLogo from "../assets/capture.png";

const Navbar = () => {
    return (
        <div className="w-full p-6 flex items-center">
            <img src={captureLogo} className="w-6 h-6" alt="Capture logo" />
            <h1 className="font-semibold text-2xl ml-1 text-[#1E1E1E]">Rekognify</h1>
        </div>
    )
}

export default Navbar


