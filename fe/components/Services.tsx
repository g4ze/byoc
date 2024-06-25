import { Service } from "@/types";
import { Heading } from "./Heading";
import { useState } from "react";
export default function Services({ services }: { services: Service[] }) {
    console.log("services from services maping component: ", services);
    const HOST_PROXY = process.env.HOST_PROXY || "localhost:80";
    const [isLoading, setIsLoading] = useState(false);
    return (
        <>
        
            {services.map((service) => (
                <section id={service.name} className="">
    <div className=" max-w-400 mx-10 text-white h-screen bg-black bg-opacity-50 items-center justify-center">
        <Heading label="Deployment Name: " className="font-bold text-sm pt-6 "/>
        <Heading
            label={service.deploymentName}
            className="font-bold text-4xl pt-2 mb-3"
        /> 
        <Heading label="Image Link: " className="font-bold text-sm pt-6 "/>
        <Heading
            label={service.image}
            className="text-2xl pt-2 mb-3"
        />
        <Heading label="DNS A" className="font-bold text-sm pt-6 "/>
        <p
            className="text-2xl pt-2 mb-3"
        ><a href={"http://"+service.loadbalancerDNS}>{service.loadbalancerDNS}</a></p>
        <Heading label="Created At: " className="font-bold text-sm pt-6 "/>
        <Heading
            label={service.createdAt.toString()}
            className="text-2xl pt-2 mb-3"
        />
        <Heading label="URL: " className="font-bold text-sm pt-6 "/>
        <Heading label={"https://"+service.slug+"."+HOST_PROXY} className="text-2xl pt-2 mb-3"/>
        <Heading label="logs:" className="font-bold text-sm pt-6 "/>
        <Heading
            label={service.logs}
            className="text-2xl pt-2 mb-3"
        />
        <button className="text-white px-3 py-2  focus:border-r-2 rounded hover:border-r-2  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 outline-none shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300" 
        onClick={async () =>{
            setIsLoading(true)

            await fetch(`http://localhost:3001/v1/delete-container/`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                },
                body: JSON.stringify({
                    "image"    :service.image,
                    "userName": service.userName,
                }),
            }
        )
        setIsLoading(false)
        window.location.reload()
    }
        }>{
            isLoading ? "Deleting..." : "Delete"
        }
    </button>
    </div>
    

                </section>
            ))}
        </>
    );
}