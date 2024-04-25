---
author: Marcus Whybrow
title: Announcing Ray Peat Rodeo
---

Hi, Marcus here.

Welcome to Ray Peat Rodeo. Ray Peat Rodeo is a web app to read, search, and discover the public works of Ray Peat. It's an [open-source codebase](https://github.com/marcuswhybrow/ray-peat-rodeo) of interviews + a custom program that reads these interviews to automatically make a website, this website.

Ray Peat has spoken widely on nutrition, but also biology, physics, science, art, philosophy, religion, history, and politics. I discovered Ray whilst rambling through critiques of the Carnivore Diet I was, then, enjoying. I'd picked up a strange book by a Mr Danny Roddy called Hair Like A Fox. Who was this Ray Peat he kept referring to?

Cut to several months later, and for my own purposes I had transcribed an interview of Ray's I found myself revisiting. I used a plain text format called Markdown which requires no special program to open.

I began developing shorthands that would eventually become Ray Peat Rodeo: "wait, who is that Ray just mentioned? Which scientific paper did he mean?" These are the questions I had. Ray Peat Rodeo is a way to answer those questions once and for all, for everyone to read. To achieve this I've hit upon a few ways to augment an interview:

**Timestamps** — First and foremost, it's useful to have links to the time in the audio or video interview for key questions and answers. For this one may use the humble square bracket surrounding a time, e.g. `[18:32]`.

**Mentions** — Next, Ray mentions a lot of things. A lot! Using double square brackets (like a wiki-link), to be distinct from timestamps one can mark names, e.g. `[[Blake, William]]`.

**Notes** — To distingish clarifying notes one uses curly brackets, e.g. `{PUFA stands for Polyunsaturated Fatty Acids}`.

**Issues** — When one can't make sense of something one can refer to a GitHub issue ID within curly brackets, e.g. `{#12}`. GitHub issues are a way of tracking, discussing and resolving issues with a coding project.

With these features, we get something like this:

<div class="relative mb-32 mt-16 w-[500px] mx-auto">
    <div class="w-[150px] h-[100px] bg-pink-100 rounded-lg absolute -rotate-1 -left-8"></div>
    <div class="w-[100px] h-[90px] bg-pink-50 rounded-lg absolute -rotate-2 -left-32 top-32"></div>
    <div class="w-[40px] h-[30px] bg-pink-50/60 rounded absolute -rotate-3 -left-24 top-16"></div>
    <div class="w-[200px] h-[150px] bg-pink-200/60 rounded-lg absolute -rotate-2 -right-32 -bottom-16"></div>
    <div class="relative top-16 text-left rounded-lg w-[500px] overflow-hidden shadow-2xl shadow-pink-500/10 font-mono text-lg text-pink-600 bg-gradient-to-bl from-pink-300 to-white -rotate-3">
        <div class="h-8">
            <div class="pt-6 pl-8">
                <div class="rounded-full bg-red-400/40 w-4 h-4 inline-block"></div>
                <div class="rounded-full bg-red-400/40 w-4 h-4 ml-[10px] inline-block"></div>
                <div class="rounded-full bg-red-400/40 w-4 h-4 ml-[10px] inline-block"></div>
            </div>
        </div>
        <div id="frontmatter" class="px-8 pt-6 pb-6">
            <p>---</p>
            <p>
                <span class="text-red-400">speakers:</span>
            </p>
            <p class="ml-4">
                <span class="text-red-400">RP:</span> Ray Peat
            </p>
            <p class="ml-4">
                <span class="text-red-400">I:</span> Interviewer
            </p>
            <p>---</p>
            <p>
                <br/>
                <span class="text-red-500">RP:</span>
                <span class="text-purple-500">[18:32]</span>
                There was an Australian study 
                <span class="text-purple-500">{#12}</span>
                around that time.
            </p>
            <p>
                <br/>
                <span class="text-red-500">I:</span>
                And who is 
                <span class="text-purple-500">[[Blake, William]]</span>?
            </p>
        </div>
    </div>
</div>

At the top you'll notice another convention I landed upon. I define a little database of whose speaking. Each speaker has a short key (such as "RP") that can be used to mark sentences throughough the interview. And each key has a partnered value (such as "Ray Peat") that let's readers know who "RP" is.

I think that looks pretty nice, but the real beauty is this: now we've standardised the format for timestamps, mentions, notes, issues, and speakers, a computer program can understand what everything in this file means. And that's the second leg upon which Ray Peat Rodeo stands, it converts that markdown interview into a web page, like this:

<div class="mt-8 mb-16 inline-block align-top text-left w-full backdrop-blur-2xl bg-gradient-to-br from-white/90 to-gray-100/30 rounded-lg shadow-2xl shadow-purple-700/20 " >
    <div class="py-4 bg-gradient-to-r from-blue-200 to-purple-300 rounded-t-lg">
        <div class="h-8 w-3/5 mx-auto bg-gradient-to-br from-white/60 to-white/50 rounded"></div>
    </div>
    <div class="px-8 pb-10 pt-2 max-w-xl mx-auto">
        <!-- Utterance -->
        <div class="font-sans ml-1 mr-16">
            <div class="text-sm mt-8 mb-4 block text-gray-400" >
                Ray Peat
            </div>
            <div class="p-8 rounded shadow text-gray-900 bg-gray-100">
                <p>
                    <span class="text-sm px-2 py-1 rounded-md bg-gray-300 hover:bg-gray-500 text-gray-50 cursor-pointer">18:32</span>
                    There was an Australian 
                    <!-- issue -->
                    <span
                        class="
                          z-10 block transition-all m-2 p-4 hover:translate-y-1 shadow-xl hover:shadow-2xl shadow-yellow-800/20 hover:shadow-yellow-600/40 rounded-md bg-gradient-to-br from-yellow-200 from-10% to-amber-200 hover:from-yellow-100 hover:from-70% hover:to-amber-200 xl:block w-2/5 mr-[-20%] float-right clear-right text-sm relative leading-5 tracking-tight
                          cursor-pointer
                        "
                    >
                        <span class="text-yellow-900 font-bold mr-0.5">
                            <img src="/assets/images/github-mark.svg" class="h-4 w-4 inline-block relative top-[-1px] mr-0.5"/> #12
                        </span>
                        <span class="text-yellow-800">Which Australian study is Ray referring to?</span>
                    </span>
                    study around that time.
                </p>
            </div>
        </div>
        <!-- Utterance -->
        <div class="font-sans ml-16 mr-1" >
            <div class="text-sm mt-8 mb-4 block text-sky-400">Interviewer</div>
            <div class="p-8 rounded shadow text-sky-900 bg-gradient-to-br from-sky-100 to-blue-200">
                <p>
                    And who is 
                    <!-- Blake mention -->
                    <span
                        hx-trigger="load"
                        hx-target="find .popup"
                        hx-get="/api/mentionable/popup/william-blake"
                        hx-swap="innerHTML"
                        hx-select=".hx-select"
                        class="relative cursor-pointer"
                        _="
                              on mouseenter
                                remove .hidden from .popup in me
                                send stopWiggling to .label in me
                              on mouseleave
                                wait for mouseenter or 500ms
                                if the result's type is not 'mouseenter'
                                  add .hidden to .popup in me
                                end
                            "
                    >
                        <span
                            class=" label font-mono font-bold tracking-normal drop-shadow-md box-decoration-clone border-b text-sky-800 hover:text-sky-900 shadow-pink-300 border-sky-800 inline-block rotate-0 transition-all "
                        >William Blake</span>?
                        <span
                            class="
                                popup
                                bg-white shadow-2xl block absolute 
                                hidden
                                z-10 
                                overflow-hidden
                                overflow-y-auto 
                                mb-4 
                                w-[400px] h-[300px]
                                left-[calc(50%-200px)]
                                top-8
                                scrollbar
                                scrollbar-track-slate-100
                                scrollbar-thumb-slate-200
                              "
                            _="on click halt the event"
                        >
                            <span class="text-center text-gray-400 block p-8">
                                loading William Blake...
                            </span>
                        </span>
                    </span>
                </p>
            </div>
        </div>
    </div>
</div>

Try hovering  over (or taping on, if on a phone/tablet) the mention of "William Blake." A popup appears. It summaries every other mention of "William Blake" in every other known interview. In this demo you can't click through, but on a real page, you can jump about between interviews this way, or see an index of every time Ray's mentioned ol' Blake.

This is the essential nurturing marriage: human readable, non proprietory, portable, plain text, markdown files that will survive _forever_; partnered with custom code that automatically derives a website from those simple files, enabling higher-order features.

I see Ray Peat Rodeo as something that _needs_ to exist, and for free, just like Ray spoke for free. Equally, everyone has bills, and I have ideas for optional paid extra for those who wish to support this kind of thing: maybe a nice PDF or eBook containing the whole thing for offline reading, something tangential like that.

My main focus is getting as many interviews into the system as possible, to that end I'm using artificial intelligence to quickly get a large quantity of low-quality transcription, with a view to improving them over time.

Thanks for your interest in the project, you can [support me on GitHub Sponsors](https://github.com/sponsors/marcuswhybrow).

Cheers,  
Marcus. 
