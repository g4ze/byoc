"use client";
import { useEffect, useState } from 'react';
import Navbar from "@/components/Navbar";
import { TabButton, ActiveTaskButton } from "@/components/TabButton";
import CreateService from '@/components/CreateService';
import Service from '@/components/Service';

export default function ServicePage() {
    const [activeService, setActiveService] = useState('create-deployment');
    const ENDPOINT="/get-services"
    const HOST_URL="http://localhost:2001"
    
    var services: Service[] = [];
    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(HOST_URL + ENDPOINT);
                const data = await response.json();
                data.forEach((service:Service) => {
                    services.push(service);
                });
                console.log("service no: "+services.length);
            } catch (error) {
                // Handle error
                console.error('Error fetching services:', error);
            }
        };

        fetchData();
    }, []);


    // const services = [
    //     { id: 'create-deployment', label: 'Create Deployment' },
    //     { id: 'service2', label: 'Service 2' },
    //     { id: 'service3', label: 'Service 3' },
    // ];

    const scrollToService = (serviceName:string) => {
        setActiveService(serviceName);
        const element = document.getElementById(serviceName);
        if (element) {
            element.scrollIntoView({ behavior: 'smooth' });
        }
    };

    return (
        <>
            <Navbar serviceName='new Deployment'/>
            <div className="grid grid-cols-11 h-screen text-gray-400">
                <div className="col-span-2 border-r-2 border-gray-300 flex justify-center pt-8">
                    <ul className="w-full">
                        {services.map((service) => (
                            <li key={service.id} className="py-2 px-4">
                                {activeService === service.id ? (
                                    <ActiveTaskButton
                                        label={service.name}
                                        onClick={() => scrollToService(service.id)}
                                    />
                                ) : (
                                    <TabButton
                                        label={service.name}
                                        onClick={() => scrollToService(service.id)}
                                    />
                                )}
                            </li>
                        ))}
                    </ul>
                </div>
                
                <div className="col-span-9 overflow-y-auto">
                    <section id="create-deployment" className="h-screen items-center justify-center">
                       <CreateService/>
                    </section>
                    <section id="service2" className="h-screen flex items-center justify-center">
                        <h2>Service 2 Section</h2>
                    </section>
                    <section id="service3" className="h-screen flex items-center justify-center">
                        <h2>Service 3 Section</h2>
                    </section>
                </div>
            </div>
        </>
    );
}