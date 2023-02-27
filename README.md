# https://raypeat.rodeo

My effort to catalogue, compile, and transcribe the public works, speeches and interviews of Ray Peat.  
[Open an issue](https://github.com/marcuswhybrow/ray-peat-rodeo/issues) if there's a Ray Peat interview I'm missing.

This repository represents a collection of markdown transcripts built by the static-site generator [11ty](https://www.11ty.dev/).

## Installation 

```
npm install
```

## Usage

```
npm start
```

## Interview Syntax

Markdown files in `./src/content/` have additional bespoke template tag shorthands for defining who's speaking, and identifying the people, books, and URLs mentioned by the speakers.

- **Ray Peat** - All paragraphs, by default, are attributed to Ray Peat. See below to attribute a paragraph to another speaker.
- **Interviewer** - Prefix a paragraph with `!MW ` to reference a speaker defined in `speakers.MW` in the frontmatter. For example: `!MW Good morning Ray`, with `speakers: { MW: Marcus Whybrow }` attributes that paragraph to `Marcus Whybrow`. Ommitting the speaker initials attributes a paragraph to most recently specified speaker above. If the frontmatter specifies one or zero speakers, all paragraphs prefixed with an `!` will be attributed to the single speaker defined, or (if zero speakers are defined) attributed to `Host`.

People, books, and URLs should be wrapped in double square brackes (`[[Text]]`) as below. Doing so feeds these links into Ray Peat Rodeo's site-wide index.

- *People* - Link to people by surrounding their full name with double square brackets `[[William Blake]]`.
- *Books* - Link to books with the title and the primary author's full name ``[[Jerusalem -by- William Blake]]``. The `-by-` separator, and exactly one author is required. Display text (see below) defaults to Book title.
- *URLs* - Link to external URLs ``[[https://www.youtube.com/watch?v=lDr71LHO0Jo]]``.
- *DOIs* - Link to scientific papers by their DOI ``[[doi:10.5860/choice.37-5129]]`` The `doi:` prefix is required.

All `[[Links]]` may optionally override the display text with the pipe sufffix `[[William Blake|a poet]]`. Hidden links (that produce no markup) are created with an empty display text string `[[William Blake|]]`. Missing links can be declared by omitting everything before the pipe `[[|text that will eventually link to something]]`.

- *Timecodes* - `[1:20:13]` or `[0:52]`. Surround a colon-separated timecode with single square brackets and Ray Peat Rodeo will generate links directly to that time using `source` URL from the frontmatter.

## Style Guide

- **Em Dashes** - Long dashses, or em dashes, (Windows ALT code 0151) when used for parenthesis contain no spaces `While I was shopping—wandering aimlessly up and down the aisles, actually—I ran into our old neighbor.

## How to Contribute

If you wish to contribute a transcription, here's the of Ray Peat interviews and articles I found [here](https://www.selftestable.com/ray-peat-stuff/sites). Sections are ordered from easiest to the hardest in terms of the work required.

I'll shortly be incorporating the following resources into the list: [Functional Performance Systems's master list](https://www.functionalps.com/blog/2011/09/12/master-list-ray-peat-phd-interviews/); Ray Peat Forum threads [Interviews](https://raypeatforum.com/community/forums/interviews.20/), [Interview Transcripts](https://raypeatforum.com/community/categories/interview-transcripts.317/), [Resources](https://raypeatforum.com/community/forums/resources.233/); [Expulsia.com/health](https://expulsia.com/health)

### Written Articles & Interviews Outside raypeat.com
- [ ] [A Biophysical Approach to Altered Consciousness](http://www.orthomolecular.org/library/jom/1975/pdf/1975-v04n03-p189.pdf)
- [ ] [A Physiological Approach to Ovarian Cancer](http://www.encognitive.com/node/3675)
- [ ] [Age-related Oxidative Changes in the Hamster Uterus](https://www.toxinless.com/age-related-oxidative-changes-in-the-hamster-uterus.pdf)
- [x] [An Interview With Dr. Raymond Peat Part I (Ten Question Interview) & II (Mind-Body Connection) - by Karen Mcc et Matt Labosco, Greg Waitt, Wayde Curran, and Mariam](http://web.archive.org/web/20160404173506/http://raypeatinsight.com/2013/06/06/organizing-the-panic-an-interview-with-dr-ray-peat/)
- [x] [An Interview With Dr. Raymond Peat: Organizing the Panic - by Karen Mcc et Wayde Curran, Eti Csiga and Tyler Derosier](http://web.archive.org/web/20160321100221/http://raypeatinsight.com/2015/01/29/raypeat-interviews-revisited)
- [ ] [An Interview With Dr. Raymond Peat by Mary Shomon](https://web.archive.org/web/20011129110010/https://www.thyroid-info.com/articles/ray-peat.htm)
- [ ] [An Interview With Dr. Raymond Peat: Negation - by Karen Mcc](https://web.archive.org/web/20200223190350/http://www.visionandacceptance.com/negation/)
- [ ] [Articles on PubMed](http://www.ncbi.nlm.nih.gov/pubmed/?term=%22Peat+R%22[Author])
- [ ] [Carbon Monoxide: Cancer Hormone?](http://www.encognitive.com/node/13878)
- [ ] [Coconut Oil and Its Virtues](http://www.naturodoc.com/library/nutrition/coconut_oil.htm)
- [ ] [Comparison of Progesterone and Estrogen](https://web.archive.org/web/20141017002820/http://www.longnaturalhealth.com/health-articles/comparison-progesterone-and-estrogen)
- [ ] [Don’t Be Conned By The Resveratrol Scam](http://doctorsaredangerous.com/articles/dont_be_conned_by_the_resveratrol_scam.htm)
- [ ] [Energy, structure, and carbon dioxide: A realistic view of the organism](http://www.functionalps.com/blog/2011/04/23/energy-structure-and-carbon-dioxide-a-realistic-view-of-the-organism/)
- [ ] [Estriol, DES, etc](http://www.encognitive.com/node/12884)
- [ ] [Estrogen and brain aging in men and women: Depression, energy, stress](http://raypeat.com/articles/articles/estrogen-age-stress.shtml)
- [ ] [Genes, Carbon Dioxide and Adaptation](http://raypeat.com/articles/articles/genes-carbon-dioxide-adaptation.shtml)
- [ ] [Hormone Balancing: Natural Treatment and Cure for Arthritis](http://www.arthritistrust.org/Articles/Hormone%20Balancing%20Natural%20Treatment%20&%20Cure%20for%20Arthritis.pdf)
- [ ] [Oral Absorption of Progesterone](https://www.toxinless.com/ray-peat-letter-to-the-editor-oral-absorption-of-progesterone.pdf)
- [ ] [Pregnenolone](https://web.archive.org/web/20161129104337/https://www.longnaturalhealth.com/health-articles/pregnenolone)
- [ ] [Progesterone](https://web.archive.org/web/20141004034741/https://www.longnaturalhealth.com/health-articles/progesterone)
- [ ] [Progesterone: Essential to Your Well-Being](http://www.tidesoflife.com/essential.htm)
- [ ] [Signs & Symptoms That Respond To Progesterone](https://web.archive.org/web/20141006182531/https://www.longnaturalhealth.com/health-articles/signs-symptoms-respond-progesterone)
- [ ] [Stress and Water](https://raypeatforum.com/community/threads/stress-and-water.1261/)
- [ ] [The Bean Syndrome](https://www.toxinless.com/ray-peat-the-bean-syndrome.pdf)
- [ ] [The Dire Effects of Estrogen Pollution](http://www.naturodoc.com/library/hormones/estrogen_pollution.htm)
- [ ] [The Generality of Adaptogens](https://www.toxinless.com/ray-peat-the-generality-of-adaptogens.pdf)
- [ ] [Thyroid](https://web.archive.org/web/20161128234400/https://www.longnaturalhealth.com/health-articles/thyroid)
- [ ] [Using Sunlight to Sustain Life](http://www.functionalps.com/blog/2012/02/27/using-sunlight-to-sustain-life/)
- [ ] [Spanish translations of some of Peat's articles](https://bloqdnotas.blogspot.com/)
- [x] [When Western Medicine Isn't Working](https://web.archive.org/web/20180123181651/https://www.beyondtheinterview.com/article/2018/when-western-medicine-isnt-workingdifferent-insights-from-a-leader-in-health)
- [ ] [On culture, government, and social class](https://medium.com/reformermag/on-culture-government-and-social-class-306dfe8af599)
- [ ] [Thyroid Deficiency and Common Health Problems](https://180degreehealth.com/thyroid-deficiency-and-common-health-problems/)
- [ ] [Negation](https://web.archive.org/web/20151123021043/http://www.visionandacceptance.com/negation)

### Audio/Video With Existing Transcripts

Ask Your Herb Doctor
- [ ] [Cancer Treatment](https://www.toxinless.com/kmud-120217-cancer-treatment.mp3) ([partial transcript](https://www.toxinless.com/kmud-120217-cancer-treatment-partial-transcript.doc))
- [ ] [Hot flashes, Night Sweats, the Relationship to Stress, Aging, PMS, Sugar Metabolism](https://www.toxinless.com/kmud-120817-hot-flashes-night-sweats-relationship-to-stress-aging-p-m-sand-sugar-metabolism.mp3) ([transcript](https://www.toxinless.com/kmud-120817-hot-flashes-night-sweats-relationship-to-stress-aging-p-m-sand-sugar-metabolism-transcript.doc))
- [ ] [Serotonin, Endotoxins, Stress](https://www.toxinless.com/kmud-110617-serotonin-endotoxins-stress.mp3) ([transcript](https://www.toxinless.com/kmud-serotonin-endotoxins-stress-110617.doc))
- [ ] [Sugar Myths I](https://www.toxinless.com/kmud-110916-sugar-myths.mp3) ([transcript](https://www.toxinless.com/kmud-110916-sugar-myths.docx))
- [ ] [Weight Gain](https://www.toxinless.com/kmud-130215-weight-gain.mp3) ([transcript](https://www.toxinless.com/kmud-130215-weight-gain-transcription.doc))
- [ ] [Suppression of Cancer Treatments](https://www.toxinless.com/polsci-010102-suppression-of-cancer.mp3) ([transcript](https://www.toxinless.com/polsci-suppression-of-cancer-treatments-transcription.pdf))

Voice of America
- [ ] [Water](https://www.toxinless.com/voiceofamerica-130909-water.mp3) ([transcript](https://www.toxinless.com/voiceofamerica-140602-water-transcription.pdf))

### Audio/Video Interviews With (Autogenerated?) Captions

Ray Peat and Budd Weiss
- [ ] [The Biology of Carbon Dioxide](https://www.youtube.com/watch?v=r6hYLtFvmw8)

### Audio/Video Interviews Without Captions

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
- [ ] [2014-01-01 Effects of Stress and Trauma on the Body](https://www.toxinless.com/eluv-140101-effects-of-stress-and-trauma-on-the-body.mp3)
- [ ] [2008-09-18 Good Fats](https://www.toxinless.com/eluv-080918-fats.mp3)


Ask Your Herb Doctor
- [ ] [2022-11](https://www.toxinless.com/kmud-221118.mp3)
- [ ] [2022-06](https://www.toxinless.com/kmud-220617.mp3)
- [ ] [2022-05](https://www.toxinless.com/kmud-220520.mp3)
- [ ] [2022-04](https://www.toxinless.com/kmud-220415.mp3)
- [ ] [2022-02](https://www.toxinless.com/kmud-220218.mp3)
- [ ] [2022-01](https://www.toxinless.com/kmud-220121.mp3)
- [ ] [2021-09](https://www.toxinless.com/kmud-210917.mp3)
- [ ] [2021-08](https://www.toxinless.com/kmud-210820.mp3)
- [ ] [2021-07](https://www.toxinless.com/kmud-210716.mp3)
- [ ] [2021-06](https://www.toxinless.com/kmud-210618.mp3)
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
- [ ] [Acidity X Alkalinity](https://www.toxinless.com/kmud-120316-acidity-x-alkalinity.mp3)
- [ ] [Aging and Energy Reversal](https://www.toxinless.com/kmud-131220-aging-and-energy-reversal.mp3)
- [ ] [Allergy](https://www.toxinless.com/kmud-160318-allergy.mp3)
- [ ] [Altitude](https://www.toxinless.com/kmud-100716-altitude.mp3)
- [ ] [Antioxidants](https://www.toxinless.com/kmud-121019-antioxidants.mp3)
- [ ] [Antioxidant Theory and the Continued War on Cancer](https://www.toxinless.com/kmud-160916-antioxidant-theory-and-continued-war-on-cancer.mp3)
- [ ] [Authoritarianism](https://www.toxinless.com/kmud-160617-authoritarianism.mp3)
- [ ] [Blood Pressure Regulation, Heart Failure, and Muscle Atrophy](https://www.toxinless.com/kmud-120720-blood-pressure-regulation-heart-failure-muscle-atrophy.mp3)
- [ ] [Bowel Endotoxin](https://www.toxinless.com/kmud-090701-bowel-endotoxin.mp3)
- [ ] [Breast Cancer](https://www.toxinless.com/kmud-150320-breast-cancer.mp3)
- [ ] [Brain "Barriers"](https://www.toxinless.com/kmud-191018-brain-barriers.mp3)
- [ ] [California SB 277 / Degradation of the Food Supply](https://www.toxinless.com/kmud-150515-california-sb-277-degradation-of-the-food-supply.mp3)
- [ ] [California Proposition 65](https://www.toxinless.com/kmud-170915-california-proposition-65.mp3)
- [ ] [Carbon Monoxide](https://www.toxinless.com/kmud-130118-carbon-monoxide.mp3)
- [ ] [Cellular Repair](https://www.toxinless.com/kmud-120615-cellular-repair.mp3)
- [ ] [Cholesterol is an Important Molecule](https://www.toxinless.com/kmud-081201-cholesterol-is-an-important-molecule.mp3)
- [ ] [Coronavirus](https://www.toxinless.com/kmud-200320-coronavirus.mp3)
- [ ] [Current Trends on Nitric Oxide](https://www.toxinless.com/kmud-151016-current-trends-nitric-oxide.mp3)
- [ ] [Continuing Research on Urea](https://www.toxinless.com/kmud-150619-continuing-research-on-urea.mp3)
- [ ] [Critical Thinking in Academia](https://www.toxinless.com/kmud-180817-critical-thinking-in-academia.mp3)
- [ ] [Dementia and Progesterone](https://www.toxinless.com/kmud-121221-dementia-progesterone.mp3)
- [ ] [Diabetes I](https://www.toxinless.com/kmud-140221-diabetes.mp3)
- [ ] [Diabetes II and How to Restore and Protect Nerves](https://www.toxinless.com/kmud-140321-how-to-restore-and-protect-nerves.mp3)
- [ ] [Diagnosis](https://www.toxinless.com/kmud-171215-diagnosis.mp3)
- [ ] [Digestion and Emotion](https://www.toxinless.com/kmud-150116-digestion-and-emotion.mp3)
- [ ] [Economics](https://www.toxinless.com/kmud-171020-economics.mp3)
- [ ] [Education and Reeducation](https://www.toxinless.com/kmud-190920-education-reeducation.mp3)
- [ ] [Endocrinology (Part 1): Parkinson's](https://www.toxinless.com/kmud-170317-endocrinology-part1-parkinsons.mp3)
- [ ] [Endocrinology (Part 2)](https://www.toxinless.com/kmud-170421-endocrinology-part2.mp3)
- [ ] [Endocrinology (Part 3)](https://www.toxinless.com/kmud-170519-endocrinology-part3.mp3)
- [ ] [Endotoxins](https://www.toxinless.com/kmud-101119-endotoxins.mp3)
- [ ] [Energetic Interactions of Ionizing and Non-ionizing Radiation](https://www.toxinless.com/kmud-121116-energetic-interactions-ionizing-and-non-ionizing-radiation.mp3)
- [ ] [Energy Production, Diabetes and Saturated Fats](https://www.toxinless.com/kmud-111118-energy-production-diabetes-saturated-fats.mp3)
- [ ] [Environmental Enrichment - Bad Science](https://www.toxinless.com/kmud-130816-environmental-enrichment.mp3)
- [ ] [Evidence Based Medicine](https://www.selftestable.com/kmud-180921-evidence-based-medicine.mp3)
- [ ] [Exploring Alternatives](https://www.toxinless.com/kmud-160520-exploring-alternatives.mp3)
- [ ] [Female Hormones / Progesterone](https://www.toxinless.com/kmud-180119-female-hormones-progesterone.mp3)
- [ ] [Field Biology](https://www.toxinless.com/kmud-140919-field-biology.mp3)
- [ ] [Food](https://www.toxinless.com/kmud-161216-food.mp3)
- [ ] [Food Additives](https://www.toxinless.com/kmud-091001-food-additives.mp3)
- [ ] [Fukujima I](https://www.toxinless.com/kmud-110318-fukujima.mp3)
- [ ] [Fukujima II, Serotonin](https://www.toxinless.com/kmud-110415-fukujima-and-serotonin.mp3)
- [ ] [Genetics X Environment](https://www.toxinless.com/kmud-120518-genetics-x-environment.mp3)
- [ ] [Hair loss, Inflammation, Osteoporosis](https://www.toxinless.com/kmud-110715-hair-loss-inflammation-osteoporosis.mp3)
- [ ] [Hashimoto's, Antibodies, Temperature and Pulse](https://www.toxinless.com/kmud-131115-hashimotos.mp3)
- [ ] [Heart I](https://www.toxinless.com/kmud-130517-heart-1.mp3)
- [ ] [Heart II](https://www.toxinless.com/kmud-130621-heart-2.mp3)
- [ ] [Heart III](https://www.toxinless.com/kmud-130719-heart-3.mp3)
- [ ] [Herbalist Sophie Lamb](https://www.toxinless.com/kmud-190719-herbalist-sophie-lamb.mp3)
- [ ] [The hormones behind inflammation](https://www.toxinless.com/kmud-211217-the-hormones-behind-inflammation.mp3)
- [ ] [Hormone Replacement Therapy](https://www.toxinless.com/kmud-130420-hormone-replacement-therapy.mp3)
- [ ] [How irradiated cells affect other living cells in human body](https://www.toxinless.com/kmud-220318-how-irradiated-cells-affect-other-living-cells-in-human-body.mp3)
- [ ] [Inflammation I](https://www.toxinless.com/kmud-110121-inflammation-1.mp3)
- [ ] [Inflammation II](https://www.toxinless.com/kmud-110218-inflammation-2.mp3)
- [ ] [Iodine, Supplement Reactions, Hormones, and More](https://www.toxinless.com/kmud-160219-iodine-supplement-reactions-hormones.mp3)
- [ ] [Language and Criticism, Estrogen (Part 1)](https://www.toxinless.com/kmud-170616-language-criticism-estrogen.mp3)
- [ ] [Language and Criticism, Estrogen (Part 2)](https://www.toxinless.com/kmud-170721-language-criticism-estrogen-part2.mp3)
- [ ] [Learned Helplessness, Nervous System and Thyroid Questionnaire](https://www.toxinless.com/kmud-130920-learned-helplessness-nervous-system-thyroid-questionaire.mp3)
- [ ] [Lipofuscin](https://www.toxinless.com/kmud-220819-lipofuscin.mp3)
- [ ] [Longevity](https://www.toxinless.com/kmud-141017-longevity.mp3)
- [ ] [Longevity and Nootropics](https://www.toxinless.com/kmud-150821-longevity-and-nootropics.mp3)
- [ ] [Managing hormones and cancer treatment with nutrition](https://www.toxinless.com/kmud-211119-managing-hormones-and-cancer-treatment-with-nutrition.mp3)
- [ ] [Medical Misinformation](https://www.toxinless.com/kmud-181019-medical-misinformation.mp3)
- [ ] [Memory, Cognition and Nutrition](https://www.toxinless.com/kmud-140516-memory-cognition-and-nutrition.mp3)
- [ ] [Milk](https://www.toxinless.com/kmud-110819-milk.mp3)
- [ ] [Mitochondria, GABA, Herbs, and More](https://www.toxinless.com/kmud-160415-mitochondria-gaba-herbs.mp3)
- [ ] [Misconceptions relating to Serotonin and Melatonin](https://www.toxinless.com/kmud-110521-misconceptions-relating-to-serotonin-and-melatonin.mp3)
- [ ] [Nitric Oxide](https://www.toxinless.com/kmud-141121-nitric-oxide.mp3)
- [ ] [Nitric Oxide, Nitrates, Nitrites, and Fluoride](https://www.toxinless.com/kmud-151218-nitric-oxide-nitrates-nitrites-fluoride.mp3)
- [ ] [On The Back of a Tiger](https://www.toxinless.com/kmud-150717-on-the-back-of-a-tiger.mp3)
- [ ] [Particles](https://www.toxinless.com/kmud-190419-particles.mp3)
- [ ] [Palpitations and Cardiac Events](https://www.toxinless.com/kmud-130315-palpitations-and-cardiac-events.mp3)
- [ ] [Phosphate and Calcium Metabolism](https://www.toxinless.com/kmud-120921-phosphate-and-calcium-metabolism.mp3)
- [ ] [Pollution](https://www.toxinless.com/kmud-190517-pollution.mp3)
- [ ] [Positive Thinking, Sleep, and Repair](https://www.toxinless.com/kmud-180615-positive-thinking-sleep-repair.mp3)
- [ ] [Postpartum Depression](https://www.toxinless.com/kmud-190621-postpartum-depression.mp3)
- [ ] [Progesterone vs Estrogen, Listener Questions (Part 1)](https://www.toxinless.com/kmud-180316-progesterone-vs-estrogen-listener-questions.mp3)
- [ ] [Progesterone vs Estrogen, Listener Questions (Part 2)](https://www.toxinless.com/kmud-180518-progesterone-vs-estrogen-listener-questions-part2.mp3)
- [ ] [The Precautionary Principle (Part 1)](https://www.toxinless.com/kmud-170120-the-precautionary-principle.mp3)
- [ ] [The Precautionary Principle (Part 2)](https://www.toxinless.com/kmud-170217-the-precautionary-principle-part2.mp3)
- [ ] [Tryptophan](https://www.toxinless.com/kmud-191115-tryptophan.mp3)
- [ ] [Questions & Answers](https://www.toxinless.com/kmud-140117-questions-and-answers.mp3)
- [ ] [Radiation](https://www.toxinless.com/kmud-101217-radiation.mp3)
- [ ] [Rheumatoid Arthritis](https://www.toxinless.com/kmud-161021-rheumatoid-arthritis.mp3)
- [ ] [Skin Cancer](https://www.selftestable.com/kmud-181116-skin-cancer.mp3)
- [ ] [Skin Cancer Part 2](https://www.selftestable.com/kmud-181221-skin-cancer-2.mp3)
- [ ] [Skin Cancer Part 3](https://www.selftestable.com/kmud-190118-skin-cancer-3.mp3)
- [ ] [Steiner Schools and Education](https://www.toxinless.com/kmud-151120-steiner-schools-and-education.mp3)
- [ ] [Sugar I](https://www.toxinless.com/kmud-100917-sugar-1.mp3)
- [ ] [Sugar II](https://www.toxinless.com/kmud-101015-sugar-2.mp3)
- [ ] [Sugar Myths II](https://www.toxinless.com/kmud-111021-sugar-myths-2.mp3)
- [ ] [The Metabolism of Cancer](https://www.toxinless.com/kmud-160715-the-metabolism-of-cancer.mp3)
- [ ] [The Ten Most Toxic Things In Our Food](https://www.toxinless.com/kmud-090901-the-ten-most-toxic-things-in-our-food.mp3)
- [ ] [Thyroid and Polyunsaturated Fatty Acids](https://www.toxinless.com/kmud-080701-thyroid-and-polyunsaturated-fatty-acids.mp3)
- [ ] [Thyroid, Metabolism and Coconut Oil](https://www.toxinless.com/kmud-080801-thyroid-metabolism-and-coconut-oil.mp3)
- [ ] [Thyroid, Polyunsaturated Fats and Oils](https://www.toxinless.com/kmud-090401-thyroid-polyunsaturated-fats-and-oils.mp3)
- [ ] [Thinking Outside the Box - New Cancer Treatments](https://www.toxinless.com/kmud-140815-thinking-outside-the-box-new-cancer-treatments.mp3)
- [ ] [Uses of Urea](https://www.toxinless.com/kmud-150220-uses-of-urea.mp3)
- [ ] [Vaccination I](https://www.toxinless.com/kmud-140620-vaccination.mp3)
- [ ] [Vaccination II](https://www.toxinless.com/kmud-140718-vaccination-2.mp3)
- [ ] [Vibrations/Frequencies](https://www.toxinless.com/kmud-200117-frequencies-vibrations.mp3)
- [ ] [Vibrations/Frequencies Part 2](https://www.toxinless.com/kmud-200221-frequencies-vibrations-2.mp3)
- [ ] [Viruses](https://www.toxinless.com/kmud-190315-viruses.mp3)
- [ ] [Vitamin D](https://www.toxinless.com/kmud-161118-vitamin-d.mp3)
- [ ] [Water Retention and Salt](https://www.toxinless.com/kmud-111216-water-retention-salt.mp3)
- [ ] [Water Quality, Atmospheric CO2, and Climate Change](https://www.toxinless.com/kmud-160115-water-quality-atmospheric-co2-climate-change.mp3)
- [ ] [You Are What You Eat (2009)](https://www.toxinless.com/kmud-090801-you-are-what-you-eat.mp3)
- [ ] [You Are What You Eat (2014)](https://www.toxinless.com/kmud-141219-you-are-what-you-eat.mp3)

Hope for Health
- [ ] [Thyroid](https://www.toxinless.com/kkvv-081031-ray-peat.mp3)

Jodellefit
- [ ] [Cortisol, Low Testosterone](https://www.toxinless.com/jf-190601-cortisol-low-testosterone.mp3)
- [ ] [Insulin Resistance, Vegans, Low Cortisol, Bone Broth, and Coconut](https://www.toxinless.com/jf-200116-insulin-resistance-vegans-low-cortisol.mp3)
- [ ] [Listener Q&A: Calories, Cortisol, Cellulite, Exercise, Ephedra & More](https://www.toxinless.com/jf-191112-listener-qa.mp3)
- [ ] [How to Fix Your Digestion & Poop](https://www.toxinless.com/jf-190910-how-to-fix-your-digestion-poop.mp3)
- [ ] [Stress and Your Health](https://www.toxinless.com/jf-190427-stress-health.mp3)

Silicon Valley Health Institute
- [ ] [Nervous System Protect & Restore](https://youtu.be/mdLHWFJI2y0)

One Radio Network - Patrick Timpone
- [ ] [Dr. Peat Answers Questions Regarding Health, Diet and Nutrition Part 1](https://www.toxinless.com/orn-140101-nutrition-1.mp3)
- [ ] [Dr. Peat Answers Questions Regarding Health, Diet and Nutrition Part 2](https://www.toxinless.com/orn-140101-nutrition-2.mp3)
- [ ] [The Goods](https://www.toxinless.com/orn-190521-the-goods.mp3)
- [ ] [Fats and Questions](https://www.toxinless.com/orn-190124-fats-and-questions.mp3)
- [ ] [Fascinating Insights Into Mr. Thyroid](https://www.toxinless.com/orn-190917-mr-thyroid.mp3)
- [ ] [Health of the Human Body](https://www.toxinless.com/orn-190319-health-of-the-human-body.mp3)
- [ ] [Menopause, Estrogen, Thyroid, Coronavirus, and Glaucoma](https://www.toxinless.com/orn-200217-menopause-estrogen-thyroid-coronavirus-glaucoma.mp3)
- [ ] [Milk](https://www.toxinless.com/orn-190718-milk.mp3)
- [ ] [Natural Healing](https://www.toxinless.com/orn-190429-natural-healing.mp3)
- [ ] [Oxygen Saturation, Lactic Acid, Thyroid, Vaccines, PUFAs](https://www.toxinless.com/orn-200120-oxygen-saturation-lactic-acid-thyroid-vaccines-pufas.mp3)
- [ ] [A Plethora of Wide-Randing Questions](https://www.toxinless.com/orn-191015-plethora-of-wide-ranging-questions.mp3)
- [ ] [Progesterone, Estrogen, Strokes, Milk, Sugars](https://www.toxinless.com/orn-191119-progesterone-estrogen-strokes-milk-sugars.mp3)
- [ ] [Thyroid, PUFAs, OJ, and Sugar](https://www.toxinless.com/orn-190219-thyroid-pufas-oj-sugar.mp3)
- [ ] [Top Contrarian on Health](https://www.toxinless.com/orn-191217-top-contrarian-on-health.mp3)
- [ ] [Vitamin D, Thyroid, Evolving Consciously](https://www.toxinless.com/orn-190820-vitamin-d-thyroid-evolving-consciously.mp3)
- [ ] [What is a virus anyway?](https://www.toxinless.com/orn-200316-what-is-a-virus-anyway.mp3)

Politics & Science
- [ ] [A Self Ordering World](https://www.toxinless.com/polsci-101110-self-ordering-world.mp3)
- [ ] [Autoimmune and Movement Disorders](https://www.toxinless.com/polsci-120518-autoimmune-and-movement.mp3)
- [ ] [Biochemical Health]()
- [ ] [Coronavirus, immunity, and vaccines (part 1)](https://www.toxinless.com/polsci-200318-coronavirus-immunity-vaccines-part1.mp3)
- [ ] [Coronavirus, immunity, and vaccines (part 2)](https://www.toxinless.com/polsci-200324-coronavirus-immunity-vaccines-part2.mp3)
- [ ] [Coronavirus, immunity, and vaccines (part 3)](https://www.toxinless.com/polsci-200331-coronavirus-immunity-vaccines-part3.mp3)
- [ ] [Digestion](https://www.toxinless.com/polsci-100426-digestion.mp3)
- [ ] [Empiricism vs Dogmatic Modeling](https://www.toxinless.com/polsci-080724-dogmatism-in-science.mp3)
- [ ] [Evolution](https://www.toxinless.com/polsci-150304-evolution.mp3)
- [ ] [Fats](https://www.toxinless.com/polsci-080721-fats.mp3)
- [ ] [Food Quality](https://www.toxinless.com/polsci-120107-food-quality.mp3)
- [ ] [Ionizing Radiation in Context Part 1](https://www.toxinless.com/polsci-090420-radiation-1.mp3)
- [ ] [Ionizing Radiation in Context Part 2](https://www.toxinless.com/polsci-090427-radiation-2.mp3)
- [ ] [Nuclear Disaster](https://www.toxinless.com/polsci-110316-nuclear-disaster.mp3)
- [ ] [Obfuscation of Radiation Science by Industry](https://www.toxinless.com/polsci-110330-obfuscation-of-radiation.mp3)
- [ ] [On The Origin of Life](https://www.toxinless.com/polsci-000102-origin-of-life.mp3)
- [ ] [Progesterone Part 1](https://www.toxinless.com/polsci-120122-progesterone-1.mp3)
- [ ] [Progesterone Part 2](https://www.toxinless.com/polsci-120129-progesterone-2.mp3)
- [ ] [Progesterone Part 3](https://www.toxinless.com/polsci-120207-progesterone-3.mp3)
- [ ] [Questions & Answers I](https://www.toxinless.com/polsci-130220-questions-and-answers.mp3)
- [ ] [Reductionist Science (5 minute excerpt)](https://www.toxinless.com/polsci-080918-reductionist-science.mp3)
- [ ] [Thyroid and Regeneration](https://www.toxinless.com/polsci-080911-thyroid-and-regeneration.mp3)
- [ ] [Two Hour Fundraiser I](https://www.toxinless.com/polsci-120222-fundraiser-1.mp3)
- [ ] [Two Hour Fundraiser II](https://www.toxinless.com/polsci-120222-fundraiser-2.mp3)
- [ ] [William Blake and Art's Relationship to Science](https://www.toxinless.com/polsci-140226-william-blake.mp3)

Rainmaking Time
- [ ] [Energy-Protective Materials](https://www.toxinless.com/rainmaking-140602-energy-protective-materials.mp3)
- [ ] [Life Supporting Substances](https://www.toxinless.com/rainmaking-110704-life-supporting-substances.mp3)

Source Nutritional Show
- [ ] [Brain and Tissue l](https://www.toxinless.com/sourcenutritional-120512-brain-and-tissue-1.mp3)
- [ ] [Brain and Tissue ll](https://www.toxinless.com/sourcenutritional-120512-brain-and-tissue-2.mp3)

World Puja
- [ ] [Foundational Hormones](https://www.toxinless.com/wp-121123-foundational-hormones.mp3)

Your Own Health And Fitness — (Note from MarshmalloW: Someone pointed out that the following interviews were probably intended to be accessed through a "library card" that helps fund the radio program. So consider [purchasing one](http://www.yourownhealthandfitness.org/?page_id=483) to support a worthy site!
- [ ] [Nutrition and the Endocrine System (February 11, 1997)](https://www.toxinless.com/yohaf-970211-nutrition-and-the-endocrine-system.mp3)
- [ ] [Thyroid/Progesterone and Diet (November 12, 1996)](https://www.toxinless.com/yohaf-961112-thyroid-progesterone-and-diet.mp3)
- [ ] [Heart, Brain, Cancer, and Hormones (July 29, 2014)](https://www.toxinless.com/yohaf-140729-heart-brain-cancer-and-hormones.mp3)