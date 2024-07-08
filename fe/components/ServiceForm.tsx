import { Heading } from './Heading';
import { Env, Service } from '@/types';
import { uniqueNamesGenerator, Config, adjectives, animals } from 'unique-names-generator';


import React, { Dispatch, SetStateAction } from 'react';
import { useState } from 'react';
const ServiceForm = ({services, setServices, setActiveService}:{services: Service[], 
    setServices: Dispatch<SetStateAction<Service[]>>, setActiveService:Dispatch<SetStateAction<string>>}) => {

  const customConfig: Config = {
    dictionaries: [adjectives, animals],
    separator: ' ',
    length: 2,
  };
    const env: Env = {
        'any':'any'
    }
  const [imageName, setImageName] = useState('');
  const [port, setPort] = useState(0);
  const [deploymentName, setDeploymentName] = useState(uniqueNamesGenerator(customConfig));
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
//   const [env, setEnv] = useState(Env);
//   const [region, setRegion] = useState('US East (Ohio)');
  const HOST_URL = process.env.NEXT_PUBLIC_BE_URL || "http://localhost:2001";
  const handleSubmit = (e:any) => {
    (async () => {
        // Handle form submission logic here
        setIsLoading(true);
        setError('');
        console.log(imageName, port);
        const resp=await fetch(HOST_URL+'/v1/deploy-container', {
            'method': 'POST',
            'headers': {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('token'),
            },
            'body': JSON.stringify({
                'image': imageName.trim(),
                'port': Number(port),
                'env': env,
                'deploymentName': deploymentName,
            }),
        });
        setIsLoading(false);
        if (resp.ok) {
            // Redirect to the service page
            console.log(resp)
            
            const data = await resp.json();
            console.log(data);
            setServices([...services, data]);
            const scrollToService = (serviceName: string) => {
                setActiveService(serviceName);
                // can do better prop drilling
                localStorage.setItem('activeService', serviceName);
                const element = document.getElementById(serviceName);
                if (element) {
                    element.scrollIntoView({ behavior: 'smooth' });
                }
            };
            scrollToService(data.name);
        } else {
            const data = await resp.json();
            setIsLoading(false);
            setError(data.message);

        }
    })();
    }


return (
    <div className=" max-w-400 mx-10 text-white ">
        <Heading label="Create Deployment: " className="font-bold text-1xl pt-6 mb-3"/>
        <input
                    type="text"
                    id="deploymentName"
                    value= {deploymentName}
                    onChange={(e) => {setDeploymentName(e.target.value)
                        console.log(deploymentName)
                    }}
                    placeholder="e.g., docker.io/library/nginx:latest"
                    className="text-4xl font-bold w-full bg-transparent text-white px-3 py-2 focus:border-r-2 rounded  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 focus:outline-none focus:shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300"/>
               
        <form>
            <div >
                <Heading label={'Image Link'} className="font-bold text-1xl py-2 text-white"/>
                <input
                    type="text"
                    id="imageName"
                    onChange={(e) => {setImageName(e.target.value)
                        console.log(imageName)
                    }}
                    placeholder="e.g., docker.io/library/nginx:latest"
                    className="text-sm w-full bg-transparent text-white px-3 py-2 focus:border-r-2 rounded  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 focus:outline-none focus:shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300"/>
                <p className="flex py-2 text-white text-sm">
                    All resources required are managed internally.
                </p>
            </div>
            <div className='flex'>
            <Heading label={'Port'} className="font-bold text-white text-1xl py-2 px-1"/>
            <input
                    type="text"
                    id="port"
                    onChange={(e) => setPort(Number(e.target.value))}
                    placeholder="e.g., 3000"
                    className="text-sm w-full bg-transparent text-white px-3 py-2 focus:border-r-2 rounded  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 focus:outline-none focus:shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300"/>
                
            </div>
            <div className="info-box py-1">
                <p>Your Free Tier project is created with a single Read/Write <br/>compute that automatically scales to zero after five minutes of inactivity.</p>
            </div>
            
            
            <div className="button-group py-2">
                    <button 
                        type="submit"
                        className="p-2 border-2 rounded mx-2 hover:bg-blue-500 hover:border-none text-white font-bold py-2 px-4 rounded"
                        disabled={isLoading}
                        onClick={handleSubmit}
                    >
                        {isLoading ? 'Creating...' : 'Create Deployment'}
                    </button>
                </div>
                {error && (
                    <div className="text-red-500 mt-2">
                        {error}
                    </div>
                )}
        </form>
    </div>
);
};

export default ServiceForm;