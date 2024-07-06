"use client";
import { useEffect, useState } from 'react';
import Navbar from "@/components/Navbar";
import { TabButton, ActiveTaskButton } from "@/components/TabButton";
import CreateService from '@/components/CreateService';
import Services from '@/components/Services';
import  {Service}  from '@/types';
import ServiceForm from '@/components/ServiceForm';

export default function ServicePage() {
    const [activeService, setActiveService] = useState('create-deployment');
    const [services, setServices] = useState<Service[]>([]);
    const ENDPOINT="/v1/get-services"
    const HOST_URL="http://localhost:2001"
    
    
    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(
                    HOST_URL + ENDPOINT,{
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': 'Bearer ' + localStorage.getItem('token')
                        }

                    }
                );
                const data = await response.json();
                console.log("data: ", data);
                data.forEach((service: Service) => {
                    console.log("service: ", service);
                    const s = [...services];
                    s.push(service);
                    setServices(s);
                });
                console.log("service no: " + services.length);
                console.log("services: ", services);
            } catch (error) {
                // Handle error
                console.error('Error fetching services:', error);
            }
        };
        fetchData();
    }, );



    const scrollToService = (serviceName:string) => {
        setActiveService(serviceName);
        const element = document.getElementById(serviceName);
        if (element) {
            element.scrollIntoView({ behavior: 'smooth' });
        }
    };

    return (
        <>
            <Navbar Name='new Deployment'/>
            <div className="grid grid-cols-11 h-screen text-gray-400">
                <div className="col-span-2 border-r-2 border-gray-300 flex justify-center pt-8">
                    <ul className="w-full">

                    <li key={"create-deployment"} className="py-2 px-4">
                                {activeService === "create-deployment" ? (
                                    <ActiveTaskButton
                                        label={"create-deployment"}
                                        onClick={() => scrollToService("create-deployment")}
                                    />
                                ) : (
                                    <TabButton
                                        label={"create-deployment"}
                                        onClick={() => scrollToService("create-deployment")}
                                    />
                                )}
                                </li>

                        {services.map((service) => (
                            
                            <li key={service.id} className="py-2 px-4">
                                {activeService === service.id ? (
                                    <ActiveTaskButton
                                        label={service.deploymentName}
                                        onClick={() => scrollToService(service.name)}
                                    />
                                ) : (
                                    <TabButton
                                        label={service.deploymentName}
                                        onClick={() => scrollToService(service.name)}
                                    />
                                )}
                            </li>
                        ))}
                    </ul>
                </div>
                
                <div className="col-span-9 overflow-y-auto">
                    <section id="create-deployment" className="h-screen items-center justify-center">
                       <ServiceForm services={services} setServices={setServices} setActiveService={setActiveService}/>
                    </section>
                    <Services services={services} setActiveService={setActiveService}/>
                </div>
            </div>
        </>
    );
}