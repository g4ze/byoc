export function Field({label, placeholder, value, onChange}:{label: string, placeholder: string, value: string, onChange: (e: any) => void}) {
    return <div className="flex flex-col">
    <label className='block text-gray-700 text-sm font-bold mt-3' >{label}</label>
    <input onChange={onChange} className='shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline' id={label} type={label} placeholder={placeholder}/>
    </div>
}