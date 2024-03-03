/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("xzzzptkywwwyxtr")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "n2zibymv",
    "name": "ingress_tls_secrets",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("xzzzptkywwwyxtr")

  // remove
  collection.schema.removeField("n2zibymv")

  return dao.saveCollection(collection)
})
