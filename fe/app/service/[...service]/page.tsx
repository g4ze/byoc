"use client"
import { usePathname } from 'next/navigation'
import { useEffect, useState } from 'react';
export default function Service() {
    // get the route
    const pathname = usePathname();
    const serviceName = pathname.split("/").pop();

    const [services, setServices] = useState([]);
    const HOST_URL = process.env.NEXT_PUBLIC_BE_URL || "http://localhost:2001";
    useEffect(() => {
        const fetchData = async () => {
            await fetch(HOST_URL+`/v1/get-services`
                ,{
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('token')
                    }

                }
            )
                .then(response => response.json())
                .then(data => setServices(data));
        };

        fetchData();
    }, [serviceName]);
    
    return(
        services
    )
}