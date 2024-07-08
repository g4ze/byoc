"use client";
import { Field } from "@/components/Field";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Heading } from "@/components/Heading";
import { Button } from "@/components/Button";

export default function Login() {
  const router = useRouter();
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");
  // URL of the login endpoint
  const loginEndpoint = "/login";
  const hosturl = process.env.NEXT_PUBLIC_BE_URL || "http://localhost:2001";
  return (
    <section className="h-screen">
      <div className="flex justify-center items-center h-screen  ">
        <div className="max-w-md w-full mx-auto bg-black bg-opacity-30">
          <div className="rounded-lg shadow-lg p-8">
            <nav className="pb-3">
              <div className="container mx-auto flex items-center justify-between">
                  Paypal
              </div>
            </nav>
            <form>

              <Heading label="Login" />
              <Field
                label="Username"
                value="text"
                placeholder="doe67"
                onChange={(e) => {
                  setUsername(e.target.value);
                }}
              />
              
              <Field
                label="Password"
                value="password"
                placeholder="********"
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
              />
              
              <div className="flex items-center justify-between">
                <Button
                  label={"Login"}
                  onClick={async () => {
                    try{
                    const response = await fetch(
                      hosturl+loginEndpoint,
                      {
                        method: "POST",
                        headers: {
                          "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                          username: username,
                          password: password,
                        }),
                      }
                    );
                    if (response.status === 200) {
                      console.log("User signed in successfully");
                      const data = await response.json();
                      localStorage.setItem("token", data.token);
                      router.push("/service");
                    } else {
                      console.log("data", await response.json());
                      alert("Invalid username or password");
                    }
                  }
                  catch(e){
                    console.log(e);
                  }
                  }}
                />
              </div>
            </form>
            <div className="text-center mt-4">
                Don&apos;t have an account? <a onClick={() => router.push("/signup")}>Sign Up</a>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
