const Footer = () => {
    const year = new Date().getFullYear()
    return (
        <div className="bg-[#1E1E1E] flex justify-center p-10">
            <p className="text-white text-sm">Copyright © {year} | All Rights Reserved</p>
        </div>
    )
}

export default Footer