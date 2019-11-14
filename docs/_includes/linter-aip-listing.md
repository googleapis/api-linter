{% comment %} When this file is included, it _must_ be sent an `aip` value.
This determines which rules are rendered.{% endcomment -%}

Rules for [AIP-{{include.aip}}](), covering
{{ page.prose_title | default: page.title | downcase }}.

<!-- prettier-ignore -->
<div class="aip-rule-listing">
  {% assign rule_pages = site.pages | where_exp: "p", "p.rule != nil" | where_exp: "p", "p.rule.aip == include.aip" | sort: "rule.name" -%}
  <table class="glue-table--datatable glue-table--stacked api-linter-rule-listing" style="width: 100%;">
    <tr>
      <th>Rule name</th>
      <th>Description</th>
    </tr>
    {% for p in rule_pages -%}
    <tr>
      <td style="vertical-align: top;">
        <a href="{{ site.url }}{{ p.url }}">
        <tt>{{ p.rule.name | last }}</tt>
        </a>
      </td>
      <td>{{ p.rule.summary }}</td>
    </tr>
    {% endfor -%}
  </table>
</div>

**Note:** Because AIPs sometimes cover topics that have some overlap, some
rules related to {{ page.prose_title | default: page.title | downcase }} may be
included in the rules for other AIPs.

[aip-{{include.aip}}]: https://aip.dev/{{ includes.aip }}
