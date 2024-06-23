export function Heading({label, className}: {label: string,className?: string}) {
  if (className) {
    return <div className={`${className}`}>
      {label}
      </div>
      }
    return <div className="font-bold text-4xl pt-6 mb-3">
      {label}
    </div>
}