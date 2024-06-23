"use client";
import { Field } from "@/components/Field";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Heading } from "@/components/Heading";
import { Button } from "@/components/Button";

export default function Signup() {
    const router = useRouter();
    const [password, setPassword] = useState("");
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    // URL of the signup login endpoint
    const signupEndpoint = "/create-user";
    const hosturl =  "http://localhost:2001";
    return (
        <section className="h-screen">
            <div className="flex justify-center items-center h-screen  bg-black bg-opacity-30">
                <div className="max-w-md w-full mx-auto ">
                    <div className="rounded-lg shadow-lg p-8">
                        <nav className="pb-3">
                            <div className="container mx-auto flex items-center justify-between">
                                Paypal
                            </div>
                        </nav>
                        <form>

                            <Heading label="Signup" />
                            <Field
                                label="Username"
                                value="text"
                                placeholder="doe67"
                                onChange={(e) => {
                                    setUsername(e.target.value);
                                }}
                            />
                            <Field
                                label="Email"
                                value="Email"
                                placeholder="********"
                                onChange={(e) => {
                                    setEmail(e.target.value);
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
                                    label={"Sign In"}
                                    onClick={async () => {
                                        console.log("signing up user: ", username, email, password)
                                        console.log("credentials: ", username, password, email)
                                        const response = await fetch(
                                            hosturl + signupEndpoint,
                                            {
                                                method: "POST",
                                                headers: {
                                                    "Content-Type": "application/json",
                                                },
                                                body: JSON.stringify({
                                                    username: username,
                                                    email: email,
                                                    password: password,
                                                    
                                                }),
                                            }
                                        );
                                        if (response.status === 200) {
                                            console.log("User signed up successfully");
                                            router.push("/login");
                                        } else {
                                            console.log("Data: ", response.text());
                                            alert("Invalid username or password");
                                        }
                                    }}
                                />
                            </div>
                        </form>
                        <div className="text-center mt-4">
                            Don't have an account? Sign Up
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
}
