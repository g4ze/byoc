import React, { useState } from 'react';
import { Heading } from './Heading';
import { Field } from './Field';
import Image from 'next/image';
import exclamationImg from '@/public/exclamation.png'

const ServiceForm = () => {
    const env: Env = {
        'any':'any'
    }
  const [imageName, setImageName] = useState('');
  const [port, setPort] = useState('0000');
//   const [env, setEnv] = useState(Env);
//   const [region, setRegion] = useState('US East (Ohio)');

  const handleSubmit = (e:any) => {
    (async () => {
        e.preventDefault();
        // Handle form submission logic here
        console.log(imageName, port);
        const resp=await fetch('http://localhost:2001/v1/deploy-container', {
            'method': 'POST',
            'headers': {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('token'),
            },
            'body': JSON.stringify({
                'imageName': imageName,
                'port': port,
                'env': env,
            }),
        });
        console.log(resp.text());
    })();
    }


return (
    <div className=" max-w-400 mx-10 text-white ">
        <Heading label="Create Deployment"/>
        <form onSubmit={handleSubmit}>
            <div >
                <Heading label={'Image Link'} className="font-bold text-1xl py-2 text-white"/>
                <input
                    type="text"
                    id="imageName"
                    value={imageName}
                    onChange={(e) => setImageName(e.target.value)}
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
                    onChange={(e) => setPort(e.target.value)}
                    placeholder="0000"
                    className="text-sm w-full bg-transparent text-white px-3 py-2 focus:border-r-2 rounded  hover:shadow-[inset_-32px_0_32px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-150 focus:outline-none focus:shadow-[inset_-32px_0_20px_-15px_rgba(255,255,255,0.2)] transition-shadow duration-300"/>
                
            </div>
            <div className="info-box py-1">
                <p>Your Free Tier project is created with a single Read/Write <br/>compute that automatically scales to zero after five minutes of inactivity.</p>
            </div>
            
            
            <div className="button-group py-2">
                <button type="submit"
                 className="p-2 border-2 rounded mx-2 hover:bg-blue-500 hover:border-none text-white font-bold py-2 px-4 rounded"
                 onClick={handleSubmit}
                 >Create Deployment
                 </button>
            </div>
        </form>
    </div>
);
};

export default ServiceForm;