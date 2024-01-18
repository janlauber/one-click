/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("a3bwaaqkbj0nfdu")

  // remove
  collection.schema.removeField("e70d5dls")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "nrwdaqpz",
    "name": "startDate",
    "type": "date",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "min": "",
      "max": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "grnqhwti",
    "name": "endDate",
    "type": "date",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": "",
      "max": ""
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("a3bwaaqkbj0nfdu")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "e70d5dls",
    "name": "active",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  // remove
  collection.schema.removeField("nrwdaqpz")

  // remove
  collection.schema.removeField("grnqhwti")

  return dao.saveCollection(collection)
})
