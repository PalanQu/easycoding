from jinja2 import Template


# hql="""
# aaa
# {% if ({{table}} != 'ods_rex_as_hammurabi_events_log_hourly') and ({{hour}} != 23 or {{table}} not in ['ods_rex_as_hammurabi_events_log_shield_hourly'] )%}
# ANALYZE TABLE 
# {% endif %}
# """

# hql="""
# aaa
# {% if {{table}} != 'ods_rex_as_hammurabi_events_log_hourly' %}
# ANALYZE TABLE 
# {% endif %}
# """
# template = Template(hql)

# s = template.render(table='test', hour='23')
# print(s)

hql= """
Hello 
{% if {{ name }}=='1' %}
{{ name }}!
{{% endif %}}
"""

template = Template(hql)
s = template.render(name='1')
print(s)
