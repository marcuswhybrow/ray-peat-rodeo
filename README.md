[raypeat.rodeo](https://raypeat.rodeo) is the open-source effort to transcribe the public works of Ray Peat. The transcripts are stored in `./documents` as markdown files. `./main.go` compiles them into static HTML. Netlify hosts, continuously deploying from the `main` branch on GitHub. Once built and hosted, Ray Peat Rodeo has the following stand-out features:

- Full text search — of all transcripts (thanks to [Pagefind](https://pagefind.app/))
- Citations — scientific papers, books, URLs and people all have direct links (see [Citations](#citations))
- Index — side-wide lookup of all [Citations](#citations)
- Timecodes — Direct links to exact moments in YouTube videos and mp3 source material. (See [Timecodes](#timecodes))
- Speakers — intelligent separation of paragraphs by who's speaking (see [Speakers](#speakers))
- Insights — top citations; transcripts at-a-glance; orderd chronologically

# Installation & Usage

```
go install             # Installs dependencies
go run main.go build   # Builds website to ./build
go run main.go dev     # Auto-builds to ./build on sourcecode changes
                       # and launches dev server with hot-reloading
go run main.go clean   # Deletes ./build and other ./lib/bin (see below)
```

[Pagefind](https://pagefind.app/), [Modd](https://github.com/cortesi/modd), and [Devd](https://github.com/cortesi/devd) are external binaries used to build the project or run the dev server. The first time you run `go run main.go build`, or `go run main.go dev`, the relevent binaries will be downloaded to `./lib/bin` and remain there for future use. `go run main.go clean` removes this directory.

# Contributing Transcripts

Add a markdown file to `./documents` with the filename format `YYYY-MM-DD-SLUG.md`, for example `2023-03-08-my-example.md`. The esssential frontmatter includes the following, where `series` is the name of the website, podcast, radio show, or indivual who is interviewing Ray. Ommit `transcription.source` if this transcription has not be published elsewhere.

```
---
title: Thyroid Function
source: https://www.youtube.com/watch?v=x6bxGBR1Sqs
series: Gary Null
transcription:
  source: https://raypeatforum.com/community/threads/ray-peat-1996-interview.4351/
  author: ilovethesea
  date: 2014-07-25
---
```

Ray Peat Rodeo has three custom markdown extensions that make transcribing easier, and let RPR understand the transcription. Namely [Speakers](#speakers), [Citations](#citations), and [Timecodes](#timecodes).

## Speakers

Use a familiar format to define who's speaking. For example, a file called `./documents/2023-03-08-example.md` might contain...

````
---
speakers:
  MW: Marcus Whybrow
---

MW: Hello Ray how are you doing?

RP: Very well, thankyou.

Did you know that speaker declarations aren't needed on every line. Only when the speaker changes.

MW: I see.
````

`RP` is a reserved shortname that always expands to `Ray Peat`. Custom speakers are defined in the markdown frontmatter as shown. If you don't know the speaker's name, use the reserved shortname `H` that expands to `Host`. Or makeup your own speaker definitions.

## Citations

Marking citations with double square brackets tells Ray Peat Rodeo to include them in the lookup index. And automatically provides the reader with a link pertenant to the citation type.

````
Citing a person such as [[William Blake]] is done with double square brackets. (outputs "William Blake")
The text [[William Blake|displayed]] can be modified whilst preserving the citation (outputs "displayed")

Book citations use the from [[Jerusalem -by- William Blake]]. (outputs "Jerusalem")
Everything to the right of "-by-" is considered the primary author's full name.
You can [[Jerusalem -by- William Blake|modify]] a book's display text too (outputs "modify")

[[doi:10.5860/choice.37-5129]] cites a scientific paper's DOI number (outputs full title of paper)
The "doi:" prefix declares everything to the right of it to be a DOI.
Again, you can [[doi:10.5860/choice.37-5129|give a pretty name too]] (outputs "give a pretty name too")

Link to URLs with [[https://raypeat.rodeo]] (outputs full title of web page)
Anything beginning with "https://" or "http://" is valid.
And it can be modified [[https://raypeat.rodeo|like any citation]] (outputs "like any citation")
If a link shouldn't appear in the lookup, index use regular markdown [link](https://google.com) syntax.

RP: Citations can be [[https://raypeat.rodeo|combined]] with speaker declarations (see start of line).
````

## Timecodes

Timecodes can be placed anywhere in a document who's source is a YouTube Video, or an mp3 file. Changes in the coversation are good places to specifiy a timecode.

````
---
source: https://www.youtube.com/watch?v=x6bxGBR1Sqs
---

[00:00:54] is a timecode.  (outputs a link to https://www.youtube.com/watch?v=x6bxGBR1Sqs#t=00h00m54s)
[00:54] or, even [54], outpus the same thing.
````

# Master List of Ray Peat Interviews

If you wish to contribute a transcription, I've compilled a list of every Ray Peat interview I'm aware of. Sections are ordered from easiest to the hardest in terms of the work required. I've checked [these sources](#sources).

## Written Articles & Interviews Outside raypeat.com
- [ ] [1975-??-?? A Biophysical Approach to Altered Consciousness](http://www.orthomolecular.org/library/jom/1975/pdf/1975-v04n03-p189.pdf)
- [ ] [2010-01-03 A Physiological Approach to Ovarian Cancer](https://web.archive.org/web/20100103100715/http://www.encognitive.com/node/3675) (wayback machine date)
- [ ] [1972-??-?? Age-related Oxidative Changes in the Hamster Uterus](https://www.toxinless.com/age-related-oxidative-changes-in-the-hamster-uterus.pdf)
- [x] [2012-07-07 Ten Question Interview (Karen Mcc, Matt Labosco, Greg Waitt, Wayde Curran, and Mariam)](https://raypeatinsight.wordpress.com/2013/06/06/raypeat-interviews-revisited/) ([See also](https://web.archive.org/web/20130612080631/http://www.dannyroddy.com/main/an-interview-with-dr-raymond-peat))
- [x] [2012-09-26 Mind-Body Connection (Karen Mcc, Matt Labosco, Greg Waitt, Wayde Curran, and Mariam)](https://raypeatinsight.wordpress.com/2013/06/06/raypeat-interviews-revisited/) ([See also](https://web.archive.org/web/20130922063825/http://www.dannyroddy.com/main/an-interview-with-dr-raymond-peat-part-ii))
- [x] [2013-03-07 Organizing the Panic (Karen Mcc, Wayde Curran, Eti Csiga and Tyler Derosier)](https://raypeatinsight.wordpress.com/2015/01/29/organizing-the-panic-an-interview-with-dr-ray-peat/)
- [ ] [2001-11-29 A Renowned Nutritional Counselor Offers His Thoughts About Thyroid Disease (Mary Shomon)](https://web.archive.org/web/20011129110010/https://www.thyroid-info.com/articles/ray-peat.htm) ([also here](https://web.archive.org/web/20010215024053/http://thyroid.about.com/health/thyroid/library/weekly/aa110800c.htm), wayback machine date)
- [ ] [2014-01-22 An Interview With Dr. Raymond Peat: Negation - by Karen Mcc](https://web.archive.org/web/20200223190350/http://www.visionandacceptance.com/negation/)
- [ ] [2010-12-15 Carbon Monoxide: Cancer Hormone?](https://web.archive.org/web/20101215070915/http://www.encognitive.com/node/13878) (wayback machine date)
- [ ] [2006-04-18 Coconut Oil and Its Virtues](https://web.archive.org/web/20060418105748/https://naturodoc.com/library/nutrition/coconut_oil.htm) (wayback machine date)
- [ ] [2014-06-20 Comparison of Progesterone and Estrogen](https://web.archive.org/web/20141017002820/http://www.longnaturalhealth.com/health-articles/comparison-progesterone-and-estrogen) (wayback machine date)
- [ ] [2009-11-03 Don’t Be Conned By The Resveratrol Scam](https://web.archive.org/web/20091103072950/https://doctorsaredangerous.com/articles/dont_be_conned_by_the_resveratrol_scam.htm) (wayback machine date)
- [ ] [2011-04-23 Energy, structure, and carbon dioxide: A realistic view of the organism](http://www.functionalps.com/blog/2011/04/23/energy-structure-and-carbon-dioxide-a-realistic-view-of-the-organism/)
- [ ] [2010-10-09 Estriol, DES, etc](https://web.archive.org/web/20101009002218/http://www.encognitive.com/node/12884)
- [ ] [1986-??-?? Hormone Balancing: Natural Treatment and Cure for Arthritis](https://web.archive.org/web/20051109135532/http://www.arthritistrust.org/Articles/Hormone%20Balancing%20Natural%20Treatment%20&%20Cure%20for%20Arthritis.pdf)
- [ ] [1986-??-?? Oral Absorption of Progesterone](https://www.toxinless.com/ray-peat-letter-to-the-editor-oral-absorption-of-progesterone.pdf)
- [ ] [2014-06-19 Pregnenolone](https://web.archive.org/web/20141010091049/https://www.longnaturalhealth.com/health-articles/pregnenolone) (wayback machine date)
- [ ] [2014-06-19 Progesterone](https://web.archive.org/web/20141004034741/https://www.longnaturalhealth.com/health-articles/progesterone) (wayback machine date)
- [ ] [1999-04-08 Progesterone: Essential to Your Well-Being](https://web.archive.org/web/20000119132821/http://www.tidesoflife.com/essential.htm)
- [ ] [2014-06-19 Signs & Symptoms That Respond To Progesterone](https://web.archive.org/web/20141006182531/https://www.longnaturalhealth.com/health-articles/signs-symptoms-respond-progesterone) (wayback machine date)
- [ ] [2004-07-05 Stress and Water](https://web.archive.org/web/20040705040123/http://members.westnet.com.au/pkolb/peat5.htm) ([Ray Peat Forum](https://raypeatforum.com/community/threads/stress-and-water.1261/)) (wayback machine date)
- [ ] [????-??-?? The Bean Syndrome](https://www.toxinless.com/ray-peat-the-bean-syndrome.pdf) (wayback machine earliest date is 2021-10-28)
- [ ] [2006-04-27 The Dire Effects of Estrogen Pollution](https://web.archive.org/web/20060427035114/http://www.naturodoc.com/library/hormones/estrogen_pollution.htm) (wayback machine date)
- [ ] [????-??-?? The Generality of Adaptogens](https://www.toxinless.com/ray-peat-the-generality-of-adaptogens.pdf) (wayback machine earliest date is 2021-10-28)
- [ ] [2014-06-09 Thyroid](https://web.archive.org/web/20141007113008/https://www.longnaturalhealth.com/health-articles/thyroid) (wayback machine date)
- [ ] [1996-06-?? Using Sunlight to Sustain Life](http://www.functionalps.com/blog/2012/02/27/using-sunlight-to-sustain-life/) (Newsletter)
- [x] [2018-01-22 When Western Medicine Isn't Working](https://web.archive.org/web/20180123181651/https://www.beyondtheinterview.com/article/2018/when-western-medicine-isnt-workingdifferent-insights-from-a-leader-in-health)
- [ ] [2016-11-28 On culture, government, and social class](https://medium.com/reformermag/on-culture-government-and-social-class-306dfe8af599) ([also](https://web.archive.org/web/20170703192300/https://reformermag.com/on-culture-government-and-social-class-306dfe8af599))
- [ ] [2013-05-06 Thyroid Deficiency and Common Health Problems](https://180degreehealth.com/thyroid-deficiency-and-common-health-problems/)
- [ ] [2014-01-22 Negation](https://web.archive.org/web/20140302004453/http://www.visionandacceptance.com/negation)
- [ ] [2005-09-19 Thyroid Information](https://web.archive.org/web/20050919235045/http://thyroid.about.com/library/weekly/aa110800c.htm) (wayback machine date)

## Audio/Video With Existing Transcripts

Ask Your Herb Doctor
- [ ] [Cancer Treatment](https://www.toxinless.com/kmud-120217-cancer-treatment.mp3) ([partial transcript](https://www.toxinless.com/kmud-120217-cancer-treatment-partial-transcript.doc))
- [ ] [Hot flashes, Night Sweats, the Relationship to Stress, Aging, PMS, Sugar Metabolism](https://www.toxinless.com/kmud-120817-hot-flashes-night-sweats-relationship-to-stress-aging-p-m-sand-sugar-metabolism.mp3) ([transcript](https://www.toxinless.com/kmud-120817-hot-flashes-night-sweats-relationship-to-stress-aging-p-m-sand-sugar-metabolism-transcript.doc))
- [ ] [Serotonin, Endotoxins, Stress](https://www.toxinless.com/kmud-110617-serotonin-endotoxins-stress.mp3) ([transcript](https://www.toxinless.com/kmud-serotonin-endotoxins-stress-110617.doc))
- [ ] [Sugar Myths I](https://www.toxinless.com/kmud-110916-sugar-myths.mp3) ([transcript](https://www.toxinless.com/kmud-110916-sugar-myths.docx))
- [ ] [Weight Gain](https://www.toxinless.com/kmud-130215-weight-gain.mp3) ([transcript](https://www.toxinless.com/kmud-130215-weight-gain-transcription.doc))

Politics & Science

- [ ] [Suppression of Cancer Treatments](https://www.toxinless.com/polsci-010102-suppression-of-cancer.mp3) ([transcript](https://www.toxinless.com/polsci-suppression-of-cancer-treatments-transcription.pdf))

Voice of America (Sharon Kleyne)
- [ ] [Water](https://www.youtube.com/watch?v=sQ8SZ7nrCJo) ([transcript](https://www.toxinless.com/voiceofamerica-140602-water-transcription.pdf), [mp3](https://www.toxinless.com/voiceofamerica-130909-water.mp3))

## Audio/Video Interviews With (Autogenerated?) Captions

Generaative Energy Podcast (Danny Roddy & Georgi Dinkov)
- [ ] [2022-06-25 #85 Protein Restriction. Lidocaine for Hair Loss? Brain Size, Intelligence & Symptom Recognition.](https://www.youtube.com/watch?v=zZCgpw6_sRA)
- [ ] [2022-04-30 #82 What's More Toxic: PUFA or Endotoxin? Euphoria. Addiction. Meaningful Work. Current Events.](https://www.youtube.com/watch?v=U2f0WOaTmu4)
- [ ] [2021-02-27 #51 Heat Shock Proteins. Antibiotic Resistance. DHT Safety. Energy and Aging.](https://www.youtube.com/watch?v=5sop_zLMZ34)

Generative Energy (Danny Roddy)
- [ ] [2018-10-01 #34 Calcium. Phosphate. Authoritarianism. Eugenics & CIA Spymaster Allen Dulles.](https://www.youtube.com/watch?v=qA9YrebN5zY)
- [ ] [2018-03-29 #33 Optimizing The Environment with Ray Peat](https://www.youtube.com/watch?v=enNGb3c1NKk)
- [ ] [2017-09-21 #32 The CIA's Mighty Wurlitzer](https://www.youtube.com/watch?v=TtTiDVWVlQc)
- [ ] [2016-07-13 #28 The Origins of Authoritarianism](https://www.youtube.com/watch?v=MBsFWsHgSCY)
- [ ] [2016-05-21 #26 Carbon Dioxide, Redox Balance, and The Ketone Body Ratio](https://www.youtube.com/watch?v=3AI46HYJ3ro)
- [ ] [2016-12-21 #31 Safe Supplements](https://www.youtube.com/watch?v=PuSfV43Quuo)

Danny Roddy
- [ ] [2016-05-04 Reliable Thyroid Products (2 minutes)](https://www.youtube.com/watch?v=VuL4daW_fXY)
- [ ] [2016-05-26 A Bioenergetic View of Ketosis in 2-Minutes](https://www.youtube.com/watch?v=H_9UOlXww3o)

Perceive Think Act Films
- [ ] [2021-02-18 On the Back of a Tiger, Day One](https://www.youtube.com/watch?v=jqhlIOt5sUw) (published dates, not recorded dates)

NPR (Gary Null)
- [x] [1996-09-02 The Thyroid](https://www.youtube.com/watch?v=5fwYK0KISTI) ([mp3](https://www.functionalps.com/blog/wp-content/uploads/2011/09/NPRraypeatinterview1996.mp3))

Bud Weiss
- [ ] [2010-10-09 The Biology of Carbon Dioxide](https://www.youtube.com/watch?v=r6hYLtFvmw8) ([video](https://www.functionalps.com/blog/wp-content/uploads/2011/09/Video-Bud-Weiss-Ray-Peat-The-Biology-of-Carbon-Dioxide-October-9-2010.mp4), [mp3](https://www.functionalps.com/blog/wp-content/uploads/2011/09/Audio-Bud-Weiss-Ray-Peat-The-Biology-of-Carbon-Dioxide-October-9-2010.mp3))

Ask Your Herb Doctor
- [ ] [2017-06-16 Language and Criticism, Estrogen (Part 1)](https://www.youtube.com/watch?v=UsdTznJr6vU) ([mp3](https://www.toxinless.com/kmud-170616-language-criticism-estrogen.mp3))
- [ ] [2017-05-19 Endocrinology (Part 3)](https://www.youtube.com/watch?v=AsuSit7exNE) ([mp3](https://www.toxinless.com/kmud-170519-endocrinology-part3.mp3))
- [ ] [2017-04-21 Endocrinology (Part 2)](https://www.youtube.com/watch?v=xJkzU-4p6MQ) ([mp3](https://www.toxinless.com/kmud-170421-endocrinology-part2.mp3))
- [ ] [2017-03-17 Endocrinology (Part 1): Parkinson's](https://www.youtube.com/watch?v=s_pawyGJ2ZU) ([mp3](https://www.toxinless.com/kmud-170317-endocrinology-part1-parkinsons.mp3))
- [ ] [2017-02-17 The Precautionary Principle (Part 2)](https://www.youtube.com/watch?v=AMswheZPO3I) ([mp3](https://www.toxinless.com/kmud-170217-the-precautionary-principle-part2.mp3))
- [ ] [2010-09-17 Sugar I](https://www.youtube.com/watch?v=Lx96YYKvA9w&t=1893s) ([mp3](https://www.toxinless.com/kmud-100917-sugar-1.mp3), [2nd mp3](https://web.archive.org/web/20190408193128/http://westernbotanicalmedicine.com/audio/2011%20Sugar%20I,%20Cholesterol,%20Obesity,%20Heart%20Disease%20Sept%202011.mp3) found [here](https://web.archive.org/web/20160406232157/https://www.westernbotanicalmedicine.com/media.html))
- [ ] [2014-04-18 (21 minute excerpt)](https://www.youtube.com/watch?v=uCAfRC6OmYQ)
- [ ] [2010-10-15 Sugar II](https://www.youtube.com/watch?v=x25EtAKpJAQ) ([mp3](https://www.toxinless.com/kmud-101015-sugar-2.mp3))
- [ ] [2010-12-17 Radiation](https://www.youtube.com/watch?v=mNOzr30pGjw) ([mp3](https://www.toxinless.com/kmud-101217-radiation.mp3))
- [ ] [2011-03-18 Fukujima I](https://www.youtube.com/watch?v=f1UisAZmlcg) ([mp3](https://www.toxinless.com/kmud-110318-fukujima.mp3))
- [ ] [2011-04-15 Fukujima II, Serotonin, Melatonin](https://www.youtube.com/watch?v=Rp4ZLLfBSxA) ([mp3](https://www.toxinless.com/kmud-110415-fukujima-and-serotonin.mp3))
- [ ] [2011-10-21 Sugar Myths II](https://www.youtube.com/watch?v=nfAMIPtW8uo) ([mp3](https://www.toxinless.com/kmud-111021-sugar-myths-2.mp3), [duplicate video?](https://www.youtube.com/watch?v=wsI8I-Iy6iQ))
- [ ] [2011-11-18 Energy Production, Diabetes and Saturated Fats](https://www.youtube.com/watch?v=0cXp2klzPcQ) ([mp3](https://www.toxinless.com/kmud-111118-energy-production-diabetes-saturated-fats.mp3))
- [ ] [2012-02-17 Cancer Treatment](https://www.youtube.com/watch?v=4RzhS1pHs8c)
- [ ] [2011-12-16 Water Retention and Salt](https://www.youtube.com/watch?v=zwmDaU4zQ8o) ([mp3](https://www.toxinless.com/kmud-111216-water-retention-salt.mp3))
- [ ] [2012-03-16 Acidity vs Alkalinity](https://www.youtube.com/watch?v=FFwxd6i8tOk) ([mp3](https://www.toxinless.com/kmud-120316-acidity-x-alkalinity.mp3))
- [ ] [2012-06-15 Cellular Repair](https://www.youtube.com/watch?v=D1KDzPKNxc8) ([mp3](https://www.toxinless.com/kmud-120615-cellular-repair.mp3))
- [ ] [2012-07-20 Blood Pressure Regulation, Heart Failure, and Muscle Atrophy](https://www.youtube.com/watch?v=A8Ce3BdRpds) ([mp3](https://www.toxinless.com/kmud-120720-blood-pressure-regulation-heart-failure-muscle-atrophy.mp3))
- [ ] [2012-05-18 Genetics vs Environment](https://www.youtube.com/watch?v=W6xTbcFAFfc&t=1947s) ([mp3](https://www.toxinless.com/kmud-120518-genetics-x-environment.mp3))
- [ ] [2012-08-17 Temperature Set Point & Regulation, Menopause and Night Sweats](https://www.youtube.com/watch?v=G-9DTsAT0Js)
- [ ] [2012-09-21 Phosphate and Calcium Metabolism](https://www.youtube.com/watch?v=-MkG0UVZ90A) ([mp3](https://www.toxinless.com/kmud-120921-phosphate-and-calcium-metabolism.mp3))
- [ ] [2012-10-19 Antioxidants](https://www.youtube.com/watch?v=JGH0YqY9CtA) ([mp3](https://www.toxinless.com/kmud-121019-antioxidants.mp3))
- [ ] [2012-11-16 Energetic Interactions of Ionizing and Non-ionizing Radiation](https://www.youtube.com/watch?v=7Z2oeCVlJJI) ([mp3](https://www.toxinless.com/kmud-121116-energetic-interactions-ionizing-and-non-ionizing-radiation.mp3))
- [ ] [2012-12-21 Dementia and Progesterone](https://www.youtube.com/watch?v=NG4pZzoYcMM) ([mp3](https://www.toxinless.com/kmud-121221-dementia-progesterone.mp3))
- [ ] [2013-01-18 Carbon Monoxide](https://www.youtube.com/watch?v=_-7_Es7-wNk) ([mp3](https://www.toxinless.com/kmud-130118-carbon-monoxide.mp3))
- [ ] [2013-02-15 Weight Gain, Foamy Urine, Fats, Light Therapy, Dreams](https://www.youtube.com/watch?v=xcHLZSBTcCc)
- [ ] [2013-03-15 Palpitations and Cardiac Events](https://www.youtube.com/watch?v=plWIkrx5ztQ) ([mp3](https://www.toxinless.com/kmud-130315-palpitations-and-cardiac-events.mp3))
- [ ] [2013-05-17 Heart I: Heart and Hormones](https://www.youtube.com/watch?v=KFxha_WinZI) ([mp3](https://www.toxinless.com/kmud-130517-heart-1.mp3))
- [ ] [2013-06-21 Heart II: PUFA, Heart Health, Statins](https://www.youtube.com/watch?v=eDjEypXh9Oc) ([mp3](https://www.toxinless.com/kmud-130621-heart-2.mp3))
- [ ] [2013-07-19 Heart III: BPA, Estrogen, Testosterone, Cancer, Tumors](https://www.youtube.com/watch?v=6Z21WRIKs_c) ([mp3](https://www.toxinless.com/kmud-130719-heart-3.mp3))
- [ ] [2013-08-16 Environmental Enrichment - Bad Science](https://www.youtube.com/watch?v=fxNxjR-aAXo) ([mp3](https://www.toxinless.com/kmud-130816-environmental-enrichment.mp3))
- [ ] [2013-09-20 Learned Helplessness, Nervous System and Thyroid Questionnaire](https://www.youtube.com/watch?v=5pXdGi8hwcM) ([mp3](https://www.toxinless.com/kmud-130920-learned-helplessness-nervous-system-thyroid-questionaire.mp3))
- [ ] [2013-11-15 Hashimoto's, Antibodies, Temperature and Pulse](https://www.youtube.com/watch?v=OVs-SlJnzs4) ([mp3](https://www.toxinless.com/kmud-131115-hashimotos.mp3))
- [ ] [2013-12-20 Aging and Energy Reversal](https://www.youtube.com/watch?v=DEdDTAyZuJU) ([mp3](https://www.toxinless.com/kmud-131220-aging-and-energy-reversal.mp3), dates disagree)
- [ ] [2014-01-17 Questions & Answers](https://www.youtube.com/watch?v=wKF6KXQcdrg) ([mp3](https://www.toxinless.com/kmud-140117-questions-and-answers.mp3))
- [ ] [2014-02-21 Diabetes I: Diabetes, Fats, Sugars, Starch Damage](https://www.youtube.com/watch?v=TwXarjDWN-g) ([mp3](https://www.toxinless.com/kmud-140221-diabetes.mp3))
- [ ] [2014-03-21 Diabetes II: How to Restore and Protect Nerves](https://www.youtube.com/watch?v=2_r6vPXp4Cw) ([mp3](https://www.toxinless.com/kmud-140321-how-to-restore-and-protect-nerves.mp3))
- [ ] [2017-01-20 The Precautionary Principle (Part 1)](https://www.youtube.com/watch?v=KdkAj1upMBs) ([mp3](https://www.toxinless.com/kmud-170120-the-precautionary-principle.mp3))
- [ ] [2016-12-16 Food](https://www.youtube.com/watch?v=KmhGEm2KNMA) ([mp3](https://www.toxinless.com/kmud-161216-food.mp3))
- [ ] [2016-11-18 Vitamin D](https://www.youtube.com/watch?v=GOu_PdIWVPc&t=1s) ([mp3](https://www.toxinless.com/kmud-161118-vitamin-d.mp3))
- [ ] [2016-10-21 Rheumatoid Arthritis](https://www.youtube.com/watch?v=qwEY0br2SEE) ([mp3](https://www.toxinless.com/kmud-161021-rheumatoid-arthritis.mp3))
- [ ] [2016-09-16 Antioxidant Theory and the Continued War on Cancer](https://www.youtube.com/watch?v=7hy5J-_oS34) ([mp3](https://www.toxinless.com/kmud-160916-antioxidant-theory-and-continued-war-on-cancer.mp3), dates disagree)
- [ ] [2016-06-17 Authoritarianism](https://www.youtube.com/watch?v=J7fCjpZyGgo) ([mp3](https://www.toxinless.com/kmud-160617-authoritarianism.mp3))
- [ ] [2016-07-15 The Metabolism of Cancer](https://www.youtube.com/watch?v=SW3EdTscTUA) ([mp3](https://www.toxinless.com/kmud-160715-the-metabolism-of-cancer.mp3))
- [ ] [2016-02-19 Iodine, Supplement Reactions, Hormones, and More](https://www.youtube.com/watch?v=huW5IRBJKiY) ([mp3](https://www.toxinless.com/kmud-160219-iodine-supplement-reactions-hormones.mp3))
- [ ] [2016-04-15 Mitochondria, GABA, Herbs, and More](https://www.youtube.com/watch?v=4mzoRk7NX8o) ([mp3](https://www.toxinless.com/kmud-160415-mitochondria-gaba-herbs.mp3))
- [ ] [2016-03-18 Allergy](https://www.youtube.com/watch?v=fME9Jb3bAD0) ([mp3](https://www.toxinless.com/kmud-160318-allergy.mp3))
- [ ] [2016-05-20 Exploring Alternatives](https://www.youtube.com/watch?v=qbNIeuojUk0) ([mp3](https://www.toxinless.com/kmud-160520-exploring-alternatives.mp3))
- [ ] [2010-07-16 Altitude and CO2](https://www.youtube.com/watch?v=Aes7DnsTuHA) ([mp3](https://www.toxinless.com/kmud-100716-altitude.mp3))
- [ ] [2009-10-01 Food Additives](https://www.youtube.com/watch?v=uxdLc6sntsQ) ([mp3](https://www.toxinless.com/kmud-091001-food-additives.mp3), dates disagree)
- [ ] [2009-09-01 The Ten Most Toxic Things In Our Food](https://www.youtube.com/watch?v=aS3hGHHjeMQ) ([mp3](https://www.toxinless.com/kmud-090901-the-ten-most-toxic-things-in-our-food.mp3), dates disagree)
- [ ] [2009-08-01 You Are What You Eat](https://www.youtube.com/watch?v=Hvuxpkm_Mt8) ([mp3](https://www.toxinless.com/kmud-090801-you-are-what-you-eat.mp3), dates disagree)
- [ ] [2009-07-01 Bowel Endotoxin](https://www.youtube.com/watch?v=xVI-g3N45Fo) ([mp3](https://www.toxinless.com/kmud-090701-bowel-endotoxin.mp3), dates disagree)
- [ ] [2009-04-01 Thyroid, Polyunsaturated Fats and Oils](https://www.youtube.com/watch?v=7OBcXWdfx0w) ([mp3](https://www.toxinless.com/kmud-090401-thyroid-polyunsaturated-fats-and-oils.mp3), dates disagree)
- [ ] [2008-07-01 Thyroid and Polyunsaturated Fatty Acids](https://www.youtube.com/watch?v=kBhuoW8zQFo) ([mp3](https://www.toxinless.com/kmud-080701-thyroid-and-polyunsaturated-fatty-acids.mp3), dates disagree)
- [ ] [2008-08-01 Thyroid, Metabolism and Coconut Oil](https://www.youtube.com/watch?v=gda19JUrTL8), ([mp3](https://www.toxinless.com/kmud-080801-thyroid-metabolism-and-coconut-oil.mp3), dates disagree)
- [ ] [2008-12-01 Cholesterol is an Important Molecule](https://www.youtube.com/watch?v=HaSjxiuqKI4) ([mp3](https://www.toxinless.com/kmud-081201-cholesterol-is-an-important-molecule.mp3), dates disagree)
- [ ] [2016-01-15 Water Quality, Atmospheric CO2, and Climate Change](https://www.youtube.com/watch?v=D5a4QAa8I-Y) ([mp3](https://www.toxinless.com/kmud-160115-water-quality-atmospheric-co2-climate-change.mp3))
- [ ] [2015-12-18 Nitric Oxide, Nitrates, Nitrites, and Fluoride](https://www.youtube.com/watch?v=RRyDs_9DDoM) ([mp3](https://www.toxinless.com/kmud-151218-nitric-oxide-nitrates-nitrites-fluoride.mp3))
- [ ] [2015-11-20 Steiner Schools and Education](https://www.youtube.com/watch?v=-kaPHCfJIZE) ([mp3](https://www.toxinless.com/kmud-151120-steiner-schools-and-education.mp3))
- [ ] [2015-10-16 Current Trends on Nitric Oxide](https://www.youtube.com/watch?v=GqeMxHvfrYc) ([mp3](https://www.toxinless.com/kmud-151016-current-trends-nitric-oxide.mp3))
- [ ] [2015-06-19 Continuing Research on Urea](https://www.youtube.com/watch?v=665aFPEc6-0) ([mp3](https://www.toxinless.com/kmud-150619-continuing-research-on-urea.mp3))
- [ ] [2015-08-21 Longevity and Nootropics](https://www.youtube.com/watch?v=w5f5k1oy9Qs) ([mp3](https://www.toxinless.com/kmud-150821-longevity-and-nootropics.mp3))
- [ ] [2015-07-17 On The Back of a Tiger](https://www.youtube.com/watch?v=kuG11gE6HgQ) ([mp3](https://www.toxinless.com/kmud-150717-on-the-back-of-a-tiger.mp3))
- [ ] [2015-01-16 Digestion and Emotion](https://www.youtube.com/watch?v=vr9rktGm1Mg) ([mp3](https://www.toxinless.com/kmud-150116-digestion-and-emotion.mp3))
- [ ] [2015-02-20 Uses of Urea](https://www.youtube.com/watch?v=Jf8RZCaKO5U) ([mp3](https://www.toxinless.com/kmud-150220-uses-of-urea.mp3))
- [ ] [2015-03-20 Breast Cancer](https://www.youtube.com/watch?v=4ijoh1b7U_E) ([mp3](https://www.toxinless.com/kmud-150320-breast-cancer.mp3))
- [ ] [2015-05-15 California SB 277 / Degradation of the Food Supply](https://www.youtube.com/watch?v=Nn8mUiSrKR8) ([mp3](https://www.toxinless.com/kmud-150515-california-sb-277-degradation-of-the-food-supply.mp3))
- [ ] [2014-12-19 You Are What You Eat](https://www.youtube.com/watch?v=CxsSVCclRyI) ([mp3](https://www.toxinless.com/kmud-141219-you-are-what-you-eat.mp3))
- [ ] [2014-11-21 Nitric Oxide](https://www.youtube.com/watch?v=D34CrlP5ZIc) ([mp3](https://www.toxinless.com/kmud-141121-nitric-oxide.mp3))
- [ ] [2014-10-17 Longevity](https://www.youtube.com/watch?v=XEgfGJ4BAGk) ([mp3](https://www.toxinless.com/kmud-141017-longevity.mp3))
- [ ] [2014-09-19 Field Biology](https://www.youtube.com/watch?v=vmzJ_X9pafA) ([mp3](https://www.toxinless.com/kmud-140919-field-biology.mp3), dates disagree)
- [ ] [2014-08-15 Thinking Outside the Box - New Cancer Treatments](https://www.youtube.com/watch?v=Z-xoj2IDWjU) ([mp3](https://www.toxinless.com/kmud-140815-thinking-outside-the-box-new-cancer-treatments.mp3))
- [ ] [2014-07-18 Vaccination II](https://www.youtube.com/watch?v=QE7XhCBK0mo) ([mp3](https://www.toxinless.com/kmud-140718-vaccination-2.mp3))
- [ ] [2014-06-20 Vaccination I](https://www.youtube.com/watch?v=FbMBoZRvnGE) ([mp3](https://www.toxinless.com/kmud-140620-vaccination.mp3))
- [ ] [2014-05-16 Memory, Cognition and Nutrition](https://www.youtube.com/watch?v=c-XQcMhAQYM) ([mp3](https://www.toxinless.com/kmud-140516-memory-cognition-and-nutrition.mp3))


Politics & Science
- [ ] [2012-02-07 Progesterone Part 3](https://www.youtube.com/watch?v=1QuHb2hT8Ho) ([mp3](https://www.toxinless.com/polsci-120207-progesterone-3.mp3))
- [ ] [2012-01-29 Progesterone Part 2](https://www.youtube.com/watch?v=sQbAjgsEvCw) ([mp3](https://www.toxinless.com/polsci-120129-progesterone-2.mp3))
- [ ] [2012-01-22 Progesterone Part 1](https://www.youtube.com/watch?v=LSY8pbvOkP8) ([mp3](https://www.toxinless.com/polsci-120122-progesterone-1.mp3))
- [ ] [2011-03-30 Obfuscation of Radiation Science by Industry](https://www.youtube.com/watch?v=pSVhiJCSVdk) ([mp3](https://www.toxinless.com/polsci-110330-obfuscation-of-radiation.mp3))
- [ ] [2011-03-16 Nuclear Disaster](https://www.youtube.com/watch?v=VdWR0I3M2zM) ([mp3](https://www.toxinless.com/polsci-110316-nuclear-disaster.mp3))
- [ ] [2000-01-02 On The Origin of Life](https://www.youtube.com/watch?v=YhyJTLe8SNc) ([mp3](https://www.toxinless.com/polsci-000102-origin-of-life.mp3))
- [ ] [2001-01-02 Supression of Cancer](https://www.youtube.com/watch?v=9NJgekVDbZo) ([also here](https://www.youtube.com/watch?v=4dhXhI--ON0))
- [ ] [2013-02-20 Questions & Answers I](https://www.youtube.com/watch?v=cTp3fDGT96s) ([mp3](https://www.toxinless.com/polsci-130220-questions-and-answers.mp3))
- [ ] [2008-09-11 Thyroid and Regeneration](https://www.youtube.com/watch?v=6mjmyFrllFI) ([mp3](https://www.toxinless.com/polsci-080911-thyroid-and-regeneration.mp3))
- [ ] [2014-02-26 William Blake and Art's Relationship to Science](https://www.youtube.com/watch?v=FDYpjrjuQjU) ([mp3](https://www.toxinless.com/polsci-140226-william-blake.mp3))
- [ ] [2012-02-22 Fundraiser](https://www.youtube.com/watch?v=TQPkMb95NdE) ([part 1 mp3](https://www.toxinless.com/polsci-120222-fundraiser-1.mp3), [part two mp3](https://www.toxinless.com/polsci-120222-fundraiser-2.mp3))

Rainmaking Time
- [ ] [2011-07-04 Life Supporting Substances](https://www.youtube.com/watch?v=DJQF-fPIwBg) ([mp3](https://www.toxinless.com/rainmaking-110704-life-supporting-substances.mp3))
- [ ] [2014-06-02 Energy-Protective Materials](https://www.youtube.com/watch?v=H88XGaC8UvE) ([mp3](https://www.toxinless.com/rainmaking-140602-energy-protective-materials.mp3))

Source Nutritional Show
- [ ] [2012-05-12 Brain and Tissue l](https://www.youtube.com/watch?v=qZiPRYEtVtw) ([mp3](https://www.toxinless.com/sourcenutritional-120512-brain-and-tissue-1.mp3))

World Puja
- [ ] [2012-11-23 Foundational Hormones](https://www.youtube.com/watch?v=0YjNyy7z4DM) ([mp3](https://www.toxinless.com/wp-121123-foundational-hormones.mp3))

Eluv
- [ ] [2014-01-01 Effects of Stress and Trauma on the Body](https://www.youtube.com/watch?v=dEd-OmoHszg) ([mp3](https://www.toxinless.com/eluv-140101-effects-of-stress-and-trauma-on-the-body.mp3))
- [ ] [2008-09-18 Good Fats](https://www.youtube.com/watch?v=Q97Bdov9V9w) ([mp3](https://www.toxinless.com/eluv-080918-fats.mp3))


## Audio/Video Interviews Without Captions

Source Nutritional Show
- [ ] [2012-05-12 Brain and Tissue ll](https://www.youtube.com/watch?v=rL282hEbYfQ) ([mp3](https://www.toxinless.com/sourcenutritional-120512-brain-and-tissue-2.mp3))

Generative Energy
- [ ] [2016-01-20 #19 Talking with Ray Peat](https://youtu.be/HPrIPVAD6dI)

Butter Living Podcast
- [ ] [2020-02-19 Fertility, Pregnancy, and Development](https://www.toxinless.com/blp-200219-fertility-pregnancy-development.mp3)
- [ ] [2019-07-22 A Casual Conversation with Ray Peat](https://www.toxinless.com/blp-190722-a-casual-conversation-with-ray-peat.mp3)

East West
- [ ] [2013-07-17 Energy and Metabolism](https://www.toxinless.com/ewh-130717-energy-and-metabolism.mp3)
- [ ] [2011-12-15 Questions & Answers II](https://www.toxinless.com/ewh-111215-q-and-a-2.mp3)
- [ ] [2011-09-29 Cholesterol and Saturated Fats](https://www.toxinless.com/ewh-110929-cholesterol-and-saturated-fats.mp3)
- [ ] [2011-08-25 Serotonin and Endotoxin](https://www.toxinless.com/ewh-110825-serotonin-and-endotoxin.mp3)
- [ ] [2011-07-12 Questions & Answers I](https://www.toxinless.com/ewh-110712-q-and-a-1.mp3)
- [ ] [2011-06-03 Milk, Calcium and Hormones](https://www.toxinless.com/ewh-110603-milk-calcium-and-hormones.mp3)
- [ ] [2011-04-27 Glycemia, Starch and Sugar in context](https://www.toxinless.com/ewh-110427-glycemia-starch-and-sugar-in-context.mp3)
- [ ] [2011-03-15 Estrogen vs Progesterone](https://www.toxinless.com/ewh-110315-estrogen-vs-progesterone.mp3)
- [ ] [2011-02-22 The Thyroid](https://www.toxinless.com/ewh-110222-the-thryoid.mp3)
- [ ] [2011-01-18 Inflammation](https://www.toxinless.com/ewh-110118-inflammation.mp3)
- [ ] [2010-11-18 The Science Behind The Dangers of Polyunsaturated Fats](https://www.toxinless.com/ewh-101118-the-science-behind-the-dangers-of-polyunsaturated-fats.mp3)

Eluv

- [ ] [2014-01-01 Effects of Stress and Trauma on the Body Part 1 (duplicate?)](https://www.functionalps.com/blog/wp-content/uploads/2011/09/wmnfFPS_000000_000001_ultra1_335FPS1.mp3)
- [ ] [2014-01-01 Effects of Stress and Trauma on the Body Part 2 (duplicate?)](https://www.functionalps.com/blog/wp-content/uploads/2011/09/wmnfFPS_000000_000002_ultra2_335FPS2.mp3)

Bud Weiss
- [ ] [2008-09-15 (20 minutes)](https://www.functionalps.com/blog/wp-content/uploads/2011/09/Audio.-Bud-Weiss-Ray-Peat.-September-15-2008.mp3)

Perceive Think Act Films
- [ ] [2021-02-18 On the Back of a Tiger, Day Two](https://www.youtube.com/watch?v=Z3yVUELD2ZA) (published dates, not recorded dates)

Ask Your Herb Doctor
- [ ] [2022-11-18](https://www.toxinless.com/kmud-221118.mp3)
- [ ] [2022-06-17](https://www.toxinless.com/kmud-220617.mp3)
- [ ] [2022-05-20](https://www.toxinless.com/kmud-220520.mp3)
- [ ] [2022-04-15](https://www.toxinless.com/kmud-220415.mp3)
- [ ] [2022-02-18](https://www.toxinless.com/kmud-220218.mp3)
- [ ] [2022-01-21](https://www.toxinless.com/kmud-220121.mp3)
- [ ] [2021-09-17](https://www.toxinless.com/kmud-210917.mp3)
- [ ] [2021-08-20](https://www.toxinless.com/kmud-210820.mp3)
- [ ] [2021-07-16](https://www.toxinless.com/kmud-210716.mp3)
- [ ] [2021-06-18](https://www.toxinless.com/kmud-210618.mp3)
- [ ] [2021-05-21](https://www.toxinless.com/kmud-210521.mp3)
- [ ] [2021-04-16](https://www.toxinless.com/kmud-210416.mp3)
- [ ] [2021-03-19](https://www.toxinless.com/kmud-210319.mp3)
- [ ] [2021-01-15](https://www.toxinless.com/kmud-210115.mp3)
- [ ] [2020-12-18](https://www.toxinless.com/kmud-201218.mp3)
- [ ] [2020-11-20](https://www.toxinless.com/kmud-201120.mp3)
- [ ] [2020-10-16](https://www.toxinless.com/kmud-201016.mp3)
- [ ] [2020-09-18](https://www.toxinless.com/kmud-200918.mp3)
- [ ] [2020-08-21](https://www.toxinless.com/kmud-200821.mp3)
- [ ] [2020-07-17](https://www.toxinless.com/kmud-200717.mp3)
- [ ] [2020-05-15](https://www.toxinless.com/kmud-200515.mp3)
- [ ] [2019-10-18 Brain "Barriers"](https://www.toxinless.com/kmud-191018-brain-barriers.mp3)
- [ ] [2017-09-15 California Proposition 65](https://www.toxinless.com/kmud-170915-california-proposition-65.mp3)
- [ ] [2020-03-20 Coronavirus](https://www.toxinless.com/kmud-200320-coronavirus.mp3)
- [ ] [2018-08-17 Critical Thinking in Academia](https://www.toxinless.com/kmud-180817-critical-thinking-in-academia.mp3)
- [ ] [2017-12-15 Diagnosis](https://www.toxinless.com/kmud-171215-diagnosis.mp3)
- [ ] [2017-10-20 Economics](https://www.toxinless.com/kmud-171020-economics.mp3)
- [ ] [2019-09-20 Education and Reeducation](https://www.toxinless.com/kmud-190920-education-reeducation.mp3)
- [ ] [2010-11-19 Endotoxins](https://www.toxinless.com/kmud-101119-endotoxins.mp3)
- [ ] [2018-09-21 Evidence Based Medicine](https://www.selftestable.com/kmud-180921-evidence-based-medicine.mp3)
- [ ] [2018-01-19 Female Hormones / Progesterone](https://www.toxinless.com/kmud-180119-female-hormones-progesterone.mp3)
- [ ] [2011-07-15 Hair loss, Inflammation, Osteoporosis](https://www.toxinless.com/kmud-110715-hair-loss-inflammation-osteoporosis.mp3)
- [ ] [2019-07-19 Herbalist Sophie Lamb](https://www.toxinless.com/kmud-190719-herbalist-sophie-lamb.mp3)
- [ ] [2021-12-17 The hormones behind inflammation](https://www.toxinless.com/kmud-211217-the-hormones-behind-inflammation.mp3)
- [ ] [2013-04-20 Hormone Replacement Therapy](https://www.toxinless.com/kmud-130420-hormone-replacement-therapy.mp3)
- [ ] [2022-03-18 How irradiated cells affect other living cells in human body](https://www.toxinless.com/kmud-220318-how-irradiated-cells-affect-other-living-cells-in-human-body.mp3)
- [ ] [2011-01-21 Inflammation I](https://www.toxinless.com/kmud-110121-inflammation-1.mp3)
- [ ] [2011-02-18 Inflammation II](https://www.toxinless.com/kmud-110218-inflammation-2.mp3)
- [ ] [2017-07-21 Language and Criticism, Estrogen (Part 2)](https://www.toxinless.com/kmud-170721-language-criticism-estrogen-part2.mp3)
- [ ] [2022-08-19 Lipofuscin](https://www.toxinless.com/kmud-220819-lipofuscin.mp3)
- [ ] [2021-11-19 Managing hormones and cancer treatment with nutrition](https://www.toxinless.com/kmud-211119-managing-hormones-and-cancer-treatment-with-nutrition.mp3)
- [ ] [2018-10-19 Medical Misinformation](https://www.toxinless.com/kmud-181019-medical-misinformation.mp3)
- [ ] [2011-08-19 Milk](https://www.toxinless.com/kmud-110819-milk.mp3)
- [ ] [2011-05-21 Misconceptions relating to Serotonin and Melatonin](https://www.toxinless.com/kmud-110521-misconceptions-relating-to-serotonin-and-melatonin.mp3)
- [ ] [2019-04-19 Particles](https://www.toxinless.com/kmud-190419-particles.mp3)
- [ ] [2019-05-17 Pollution](https://www.toxinless.com/kmud-190517-pollution.mp3)
- [ ] [2018-06-15 Positive Thinking, Sleep, and Repair](https://www.toxinless.com/kmud-180615-positive-thinking-sleep-repair.mp3)
- [ ] [2019-06-21 Postpartum Depression](https://www.toxinless.com/kmud-190621-postpartum-depression.mp3)
- [ ] [2018-03-16 Progesterone vs Estrogen, Listener Questions (Part 1)](https://www.toxinless.com/kmud-180316-progesterone-vs-estrogen-listener-questions.mp3)
- [ ] [2018-05-18 Progesterone vs Estrogen, Listener Questions (Part 2)](https://www.toxinless.com/kmud-180518-progesterone-vs-estrogen-listener-questions-part2.mp3)
- [ ] [2019-11-15 Tryptophan](https://www.toxinless.com/kmud-191115-tryptophan.mp3)
- [ ] [2018-11-16 Skin Cancer](https://www.selftestable.com/kmud-181116-skin-cancer.mp3)
- [ ] [2018-12-21 Skin Cancer Part 2](https://www.selftestable.com/kmud-181221-skin-cancer-2.mp3)
- [ ] [2019-01-18 Skin Cancer Part 3](https://www.selftestable.com/kmud-190118-skin-cancer-3.mp3)
- [ ] [2020-01-17 Vibrations/Frequencies](https://www.toxinless.com/kmud-200117-frequencies-vibrations.mp3)
- [ ] [2020-02-21 Vibrations/Frequencies Part 2](https://www.toxinless.com/kmud-200221-frequencies-vibrations-2.mp3)
- [ ] [2019-03-15 Viruses](https://www.toxinless.com/kmud-190315-viruses.mp3)

Hope for Health
- [ ] [2008-10-31 Thyroid](https://www.toxinless.com/kkvv-081031-ray-peat.mp3)

Jodellefit
- [ ] [2019-06-01 Cortisol, Low Testosterone](https://www.toxinless.com/jf-190601-cortisol-low-testosterone.mp3)
- [ ] [2020-01-16 Insulin Resistance, Vegans, Low Cortisol, Bone Broth, and Coconut](https://www.toxinless.com/jf-200116-insulin-resistance-vegans-low-cortisol.mp3)
- [ ] [2019-11-12 Listener Q&A: Calories, Cortisol, Cellulite, Exercise, Ephedra & More](https://www.toxinless.com/jf-191112-listener-qa.mp3)
- [ ] [2019-09-10 How to Fix Your Digestion & Poop](https://www.toxinless.com/jf-190910-how-to-fix-your-digestion-poop.mp3)
- [ ] [2019-04-27 Stress and Your Health](https://www.toxinless.com/jf-190427-stress-health.mp3)

Silicon Valley Health Institute
- [ ] [2005-10-?? Nervous System Protect & Restore](https://youtu.be/mdLHWFJI2y0)

One Radio Network - Patrick Timpone
- [ ] [2014-01-01 Dr. Peat Answers Questions Regarding Health, Diet and Nutrition Part 1](https://www.toxinless.com/orn-140101-nutrition-1.mp3)
- [ ] [2014-01-01 Dr. Peat Answers Questions Regarding Health, Diet and Nutrition Part 2](https://www.toxinless.com/orn-140101-nutrition-2.mp3)
- [ ] [2019-05-21 The Goods](https://www.toxinless.com/orn-190521-the-goods.mp3)
- [ ] [2019-01-24 Fats and Questions](https://www.toxinless.com/orn-190124-fats-and-questions.mp3)
- [ ] [2019-09-17 Fascinating Insights Into Mr. Thyroid](https://www.toxinless.com/orn-190917-mr-thyroid.mp3)
- [ ] [2019-03-19 Health of the Human Body](https://www.toxinless.com/orn-190319-health-of-the-human-body.mp3)
- [ ] [2020-02-17 Menopause, Estrogen, Thyroid, Coronavirus, and Glaucoma](https://www.toxinless.com/orn-200217-menopause-estrogen-thyroid-coronavirus-glaucoma.mp3)
- [ ] [2019-07-18 Milk](https://www.toxinless.com/orn-190718-milk.mp3)
- [ ] [2019-04-29 Natural Healing](https://www.toxinless.com/orn-190429-natural-healing.mp3)
- [ ] [2020-01-20 Oxygen Saturation, Lactic Acid, Thyroid, Vaccines, PUFAs](https://www.toxinless.com/orn-200120-oxygen-saturation-lactic-acid-thyroid-vaccines-pufas.mp3)
- [ ] [2019-10-15 A Plethora of Wide-Randing Questions](https://www.toxinless.com/orn-191015-plethora-of-wide-ranging-questions.mp3)
- [ ] [2019-11-19 Progesterone, Estrogen, Strokes, Milk, Sugars](https://www.toxinless.com/orn-191119-progesterone-estrogen-strokes-milk-sugars.mp3)
- [ ] [2019-02-19 Thyroid, PUFAs, OJ, and Sugar](https://www.toxinless.com/orn-190219-thyroid-pufas-oj-sugar.mp3)
- [ ] [2019-12-17 Top Contrarian on Health](https://www.toxinless.com/orn-191217-top-contrarian-on-health.mp3)
- [ ] [2019-08-20 Vitamin D, Thyroid, Evolving Consciously](https://www.toxinless.com/orn-190820-vitamin-d-thyroid-evolving-consciously.mp3)
- [ ] [2020-03-16 What is a virus anyway?](https://www.toxinless.com/orn-200316-what-is-a-virus-anyway.mp3)

Politics & Science
- [ ] [2010-11-10 A Self Ordering World](https://www.toxinless.com/polsci-101110-self-ordering-world.mp3)
- [ ] [2012-05-18 Autoimmune and Movement Disorders](https://www.toxinless.com/polsci-120518-autoimmune-and-movement.mp3)
- [ ] [2015-03-11 Biochemical Health](https://www.toxinless.com/polsci-150311-biochemical-health.mp3)
- [ ] [2020-03-18 Coronavirus, immunity, and vaccines (part 1)](https://www.toxinless.com/polsci-200318-coronavirus-immunity-vaccines-part1.mp3)
- [ ] [2020-03-24 Coronavirus, immunity, and vaccines (part 2)](https://www.toxinless.com/polsci-200324-coronavirus-immunity-vaccines-part2.mp3)
- [ ] [2020-03-31 Coronavirus, immunity, and vaccines (part 3)](https://www.toxinless.com/polsci-200331-coronavirus-immunity-vaccines-part3.mp3)
- [ ] [2010-04-26 Digestion](https://www.toxinless.com/polsci-100426-digestion.mp3)
- [ ] [2008-07-24 Empiricism vs Dogmatic Modeling](https://www.toxinless.com/polsci-080724-dogmatism-in-science.mp3)
- [ ] [2015-03-04 Evolution](https://www.toxinless.com/polsci-150304-evolution.mp3)
- [ ] [2008-07-21 Fats](https://www.toxinless.com/polsci-080721-fats.mp3)
- [ ] [2012-01-07 Food Quality](https://www.toxinless.com/polsci-120107-food-quality.mp3)
- [ ] [2009-04-20 Ionizing Radiation in Context Part 1](https://www.toxinless.com/polsci-090420-radiation-1.mp3)
- [ ] [2009-04-27 Ionizing Radiation in Context Part 2](https://www.toxinless.com/polsci-090427-radiation-2.mp3)
- [ ] [2008-09-18 Reductionist Science (5 minute excerpt)](https://www.toxinless.com/polsci-080918-reductionist-science.mp3)
- [ ] [????-??-?? Machine Scientist (5 minutes)](https://www.functionalps.com/blog/wp-content/uploads/2011/09/Machinist-Scientists.mp3)

Your Own Health And Fitness — (Note from MarshmalloW: Someone pointed out that the following interviews were probably intended to be accessed through a "library card" that helps fund the radio program. So consider [purchasing one](http://www.yourownhealthandfitness.org/?page_id=483) to support a worthy site!
- [ ] [1997-02-11 Nutrition and the Endocrine System (February 11, 1997)](https://www.toxinless.com/yohaf-970211-nutrition-and-the-endocrine-system.mp3)
- [ ] [1996-11-12 Thyroid/Progesterone and Diet (November 12, 1996)](https://www.toxinless.com/yohaf-961112-thyroid-progesterone-and-diet.mp3)
- [ ] [2014-07-29 Heart, Brain, Cancer, and Hormones (July 29, 2014)](https://www.toxinless.com/yohaf-140729-heart-brain-cancer-and-hormones.mp3)

KWAI 1080 AAM
- [ ] [2012-05-05 Interview 1](https://www.functionalps.com/blog/wp-content/uploads/2019/01/Ray-Peat-5.5.12-edited-version.mp3)
- [ ] [2012-05-12 Interview 2](https://www.functionalps.com/blog/wp-content/uploads/2019/01/Ray-Peat-5.12.12-edited-version.mp3)

# Sources

The above list incorporates (or will eventually) the following sources. [Open an issue](https://github.com/marcuswhybrow/ray-peat-rodeo/issues) to suggest more.

- [x] [MarshmalloW's Maaster List](https://www.selftestable.com/ray-peat-stuff/sites) (as of 2023-02-28)
- [x] [Functional Performance Systems's Master List](https://www.functionalps.com/blog/2011/09/12/) (as of 2023-02-28)master-list-ray-peat-phd-interviews/)
- [ ] [Ray Peat Forums -> Interviews](https://raypeatforum.com/community/forums/interviews.20/)
- [ ] [Ray Peat Forums -> Interview Transcripts](https://raypeatforum.com/community/categories/interview-transcripts.317/)
- [ ] [Ray Peat Forums -> Resources](https://raypeatforum.com/community/forums/resources.233/)
- [ ] [Expulsia.com/health](https://expulsia.com/health)
- [x] [Ray Peat Clips](https://www.youtube.com/channel/UCh4kMDfEon-IAlQcbGym9UQ/videos) (as of 2023-02-28)
- [ ] [Western Botanical Medicine List](https://web.archive.org/web/20160406232157/https://www.westernbotanicalmedicine.com/media.html) (looks like a duplicate)
- [ ] [Chadnet](https://wiki.chadnet.org/ray-peat)
- [ ] [Ray Peat Interview Resources](https://github.com/Ray-Peat/interview/wiki)

# Non-Pertenant Sources

Currently focusing on audio/video and written interviews and guest articles. These source currently beyond the scope of this project. Subject to change.

- [ ] [raypeat.com](https://raypeat.com)
- [ ] Ray Peat's Newsletters
- [ ] [Articles on PubMed](http://www.ncbi.nlm.nih.gov/pubmed/?term=%22Peat+R%22[Author])
- [ ] [Spanish translations of some of Peat's articles](https://bloqdnotas.blogspot.com/)