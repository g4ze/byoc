import { Service } from "@/types";
export default function Services({ services }: { services: Service[] }) {
    console.log("services from services maping component: ", services);
    return (
        <>
        
            {services.map((service) => (
                <section id={service.name} className="h-screen flex items-center justify-center">
                    <h2>{service.deploymentName}</h2>
                    hi
                    {/* Add the content for each service section here */}

                </section>
            ))}
        </>
    );
}