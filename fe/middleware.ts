import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'


export async function middleware(request: NextRequest) {
  const url = request.nextUrl.clone()
  const BE_HOST= process.env.BE_HOST||"http://localhost:2001"
    console.info("path: ", url.pathname)
    const paths= url.pathname.split("/")
    if (paths[1]==='d'){
        // get DNS info from DB using path slug
        try{
            const apiUrl = `${BE_HOST}/get-lbdns?slug=${paths[2]}`;
            console.log("apiUrl: ", apiUrl)
            const response = await fetch(apiUrl);
      if (response.ok) {
        const data = await response.json();
        url.hostname = data.lbDns;
        url.pathname = paths.slice(3).join('/');
        url.port = ''
        console.info(`Redirecting to: ${url.toString()}`);
      } else {
        console.error(`Error getting LB DNS: ${response.statusText}`);
      }
        }
        catch(e:any){
            console.error("Error getting DNS: "+e.message)
        }
        
    }

    
    const response = NextResponse.rewrite(url)
    
    // Add CORS headers if needed
    response.headers.set('Access-Control-Allow-Origin', '*')
    response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
    response.headers.set('Access-Control-Allow-Headers', 'Content-Type, Authorization')
    
    return response
  }


