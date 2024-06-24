export function TabButton({ label, onClick }:{label:string, onClick:()=>void}) {

    return (
        <button 
        className="text-sm w-full align-left bg-transparent text-gray-400 px-3 py-2  focus:border-r-2 rounded  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 focus:outline-none focus:shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300" 
        onClick={onClick}>
    {label}
</button>
    );
}
export function ActiveTaskButton({ label, onClick }:{label:string, onClick:()=>void}) {
    
        return <button 
        className="w-full text-white px-3 py-2  border-r-2 rounded  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 outline-none shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300"        onClick={onClick}>
            {label}
        </button>
    }

