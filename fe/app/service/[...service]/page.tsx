"use client"
import { usePathname } from 'next/navigation'
import { useEffect, useState } from 'react';
export default function Service() {
    // get the route
    const pathname = usePathname();
    const serviceName = pathname.split("/").pop();

    const [services, setServices] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            await fetch(`http://localhost:3001/v1/get-services`)
                .then(response => response.json())
                .then(data => setServices(data));
        };

        fetchData();
    }, [serviceName]);
    
    return(
        services
    )
}