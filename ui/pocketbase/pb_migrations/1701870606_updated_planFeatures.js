/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("m9zfrl1ndxsejxd")

  // remove
  collection.schema.removeField("qcxoaiao")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "joavjbxi",
    "name": "plan",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "s7nhljkrmzzu8y6",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("m9zfrl1ndxsejxd")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "qcxoaiao",
    "name": "name",
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

  // remove
  collection.schema.removeField("joavjbxi")

  return dao.saveCollection(collection)
})
