---
title: Hello World
layout: base.njk
---

Hi people things are going well.

{% for item in collections.all %}
- <a href="{{ item.url }}">{{ item.data.title }}</a>
{% endfor %}