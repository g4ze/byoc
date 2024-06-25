import { NextApiRequest, NextApiResponse } from 'next';
import { config } from 'dotenv';
import { Pool } from 'pg';

config({ path: '../../../.env.postgres' });

const pool = new Pool({
  host: 'localhost',
  port: 5432,
  user: process.env.POSTGRES_USER||"postgres",
  password: process.env.POSTGRES_PASSWORD||"Welcome",
  database: process.env.POSTGRES_DB||"postgres",
  ssl: false,
});

export async function GET(req: NextApiRequest): Promise<Response> {
  var slug=req.url?.toString().split("slug")[1]
  slug=slug?.split("=")[1]
  console.log("slug: ", slug)
    console.log("pool: ", pool)
  

  try {
    const query = 'SELECT "loadbalancerDNS" FROM "Service" WHERE "slug" = $1';
    const result = await pool.query(query, [slug]);

    if (result.rows.length === 0) {
      return new Response(JSON.stringify({ error: 'No service URL found for slug' }),
        { status: 404, headers: { 'Content-Type': 'application/json' } });
    }
    console.log("result dns of getlbdns: ", result)
    return new Response(JSON.stringify( { lbDns: result.rows[0].loadbalancerDNS }),
        { status: 200, headers: { 'Content-Type': 'application/json' } });
    
  } catch (error) {
    console.error('Error querying database:', error);
    return new Response(JSON.stringify({ error: 'Internal server error'  }),
    { status: 500, headers: { 'Content-Type': 'application/json' } });

  }
}