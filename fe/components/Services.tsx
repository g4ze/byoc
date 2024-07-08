import { Service } from "@/types";
import { Heading } from "./Heading";
import { Dispatch, SetStateAction, useState } from "react";
export default function Services({ services, setActiveService }: { services: Service[], setActiveService: Dispatch<SetStateAction<string>> }) {
    console.log("services from services maping component: ", services);
    const HOST_PROXY = process.env.HOST_PROXY || "localhost:3000";
    const HOST_URL = process.env.NEXT_PUBLIC_BE_URL || "http://localhost:2001";
    const [isLoading, setIsLoading] = useState(false);
    const [isDeleted, setIsDeleted] = useState(false);
    return (
        <>

            {services.map((service) => (
                service.createdAt && localStorage.setItem(service.slug, service.createdAt.toString()),
                <section id={service.name} className="">
                    <div className=" max-w-400 mx-10 text-white h-screen bg-black bg-opacity-50 items-center justify-center">
                        <Heading label="Deployment Name: " className="font-bold text-sm pt-6 " />
                        <Heading
                            label={service.deploymentName}
                            className="font-bold text-4xl pt-2 mb-3"
                        />
                        <Heading label="Image Link: " className="font-bold text-sm pt-6 " />
                        <Heading
                            label={service.image}
                            className="text-2xl pt-2 mb-3"
                        />
                        <Heading label="DNS A" className="font-bold text-sm pt-6 " />
                        <p
                            className="text-2xl pt-2 mb-3"
                        ><a href={"http://" + service.loadbalancerDNS}>{service.loadbalancerDNS}</a></p>
                        <Heading label="Created At: " className="font-bold text-sm pt-6 " />
                        <Heading
                            label={service.createdAt?.toString()||localStorage.getItem(service.slug)?.toString()||"recently"}
                            className="text-2xl pt-2 mb-3"
                        />
                        <Heading label="URL: " className="font-bold text-sm pt-6 " />
                        <a href={"https://" +HOST_PROXY+"/d/"+ service.slug.replace(" ","-")}><Heading label={"https://" +HOST_PROXY+"/d/"+ service.slug.replace(" ","-")} className="text-2xl pt-2 mb-3" />
                        </a>
                        <Heading label="logs:" className="font-bold text-sm pt-6 " />
                        <Heading
                            label={service.logs}
                            className="text-2xl pt-2 mb-3"
                        />
                        <button className="text-white px-3 py-2  focus:border-r-2 rounded hover:border-r-2  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 outline-none shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300"
                            onClick={async () => {
                                setIsLoading(true)

                                await fetch(HOST_URL+`/v1/delete-container`, {
                                    method: "DELETE",
                                    headers: {
                                        "Content-Type": "application/json",
                                        "Authorization": "Bearer " + localStorage.getItem("token"),
                                    },
                                    body: JSON.stringify({
                                        "image": service.image,
                                        "userName": service.userName,
                                    }),
                                }
                                )
                                setIsLoading(false)
                                setIsDeleted(true)
                                // better prop drilling can be done
                                const scrollToService = (serviceName: string) => {
                                    setActiveService(serviceName);
                                    // can do better prop drilling
                                    localStorage.setItem('activeService', serviceName);
                                    const element = document.getElementById(serviceName);
                                    if (element) {
                                        element.scrollIntoView({ behavior: 'smooth' });
                                    }
                                };
                                scrollToService("create-deployment")
                                window.location.reload()
                            }
                            }>{
                                isLoading ? "Deleting..." : isDeleted? "Deleted" : "Delete"
                            }
                        </button>
                    </div>


                </section>
            ))}
        </>
    );
}